package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	HTTP11           = "HTTP/1.1"
	HTTP10           = "HTTP/1.0"
	DocRoot          = "./www.sjsu.edu"
	DefaultPort      = "8080"
	DefaultIndexFile = "index.html"
)

var activeConnections int64

func main() {
	documentRoot := flag.String("document_root", DocRoot, "The document root to serve files from")
	port := flag.String("port", DefaultPort, "Port to listen on")
	flag.Parse()

	// Bind a port to server.
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Printf("Error starting server on port %s: %v\n", *port, err)
		return
	}
	defer listener.Close()

	// Just so I know if I'm not querying the wrong folder.
	fmt.Printf("Serving files from %s on port %s...\n", *documentRoot, *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		atomic.AddInt64(&activeConnections, 1)
		fmt.Println("The number of connections", activeConnections)
		// Spawn a new goroutine for handling each connection
		go handleConnection(conn, *documentRoot)
	}
}

// Just use the number of active connections for dynamic timeout.
func getDynamicTimeout() time.Duration {
	switch {
	case activeConnections > 500:
		return 3 * time.Second
	case activeConnections > 100:
		return 5 * time.Second
	case activeConnections > 10:
		return 10 * time.Second
	default:
		return 15 * time.Second
	}
}

// Handles new connections and tries to reuse connections if possible
func handleConnection(conn net.Conn, documentRoot string) {
	defer conn.Close()
	defer func() {
		atomic.AddInt64(&activeConnections, -1)
	}()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Handling connection for client: %s\n", clientAddr)

	reader := bufio.NewReader(conn)

	// Infinite for loop to keep receiving requests on the same tcp connection.
	for {
		conn.SetReadDeadline(time.Now().Add(getDynamicTimeout()))
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s closed the connection (EOF)\n", clientAddr)
			} else {
				handleReadError(err, conn)
			}
			return
		}

		// Read Headers
		headers := make(map[string]string)
		for {
			headerLine, err := reader.ReadString('\n')
			if err != nil || headerLine == "\r\n" {
				break
			}

			headerParts := strings.SplitN(headerLine, ":", 2)
			if len(headerParts) == 2 {
				headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
			}
		}

		if err := processRequest(conn, requestLine, headers, documentRoot); err != nil {
			fmt.Println("Error processing request:", err)
			return
		}

		// Determine if the connection should be kept alive
		protocol := strings.Split(requestLine, " ")[2]
		if !isKeepAlive(headers, protocol) {
			fmt.Printf("Closing connection for client: %s\n", clientAddr)
			return
		}
	}
}

func handleReadError(err error, conn net.Conn) {
	if os.IsTimeout(err) {
		fmt.Println("Read timeout occurred:", err)
		conn.Close()
	} else if err == io.EOF {
		fmt.Println("Client closed the connection (EOF)")
		conn.Close()
	} else {
		fmt.Printf("Error reading from connection: %v\n", err)
		conn.Close()
	}

}

func processRequest(conn net.Conn, requestLine string, headers map[string]string, documentRoot string) error {
	tokens := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(tokens) != 3 || tokens[0] != "GET" {
		sendError(conn, "400 Bad Request", HTTP10)
		return errors.New("bad request")
	}

	path := tokens[1]
	if path == "/" {
		path = "/" + DefaultIndexFile
	}

	protocol := tokens[2]
	if protocol != HTTP10 && protocol != HTTP11 {
		sendError(conn, "505 HTTP Version Not Supported", protocol)
		return errors.New("unsupported HTTP version")
	}

	// Check for the Host header as it's required by default for http/1.1.
	if protocol == HTTP11 {
		if _, ok := headers["Host"]; !ok {
			sendError(conn, "400 Bad Request", protocol)
			return errors.New("missing host header")
		}
	}

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("[%s] Client %s requested %s %s\n", time.Now().Format(time.RFC1123), clientAddr, tokens[0], tokens[1])

	keepAlive := isKeepAlive(headers, protocol)

	filePath := filepath.Join(documentRoot, filepath.Clean(path))
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(documentRoot)) {
		sendError(conn, "403 Forbidden", protocol)
		return errors.New("forbidden access")
	}

	return serveFile(conn, filePath, keepAlive, protocol)
}

// Checks if file exists, sets the success headers and returns an error if any.
func serveFile(conn net.Conn, path string, keepAlive bool, protocol string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(conn, "404 Not Found", protocol)
		} else if os.IsPermission(err) {
			sendError(conn, "403 Forbidden", protocol)
		} else {
			sendError(conn, "500 Internal Server Error", protocol)
		}
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		sendError(conn, "404 Not Found", protocol)
		return err
	}

	contentType := getContentType(filepath.Ext(path))
	dateHeader := time.Now().UTC().Format(time.RFC1123)
	headers := fmt.Sprintf("%s 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\nDate: %s\r\n",
		protocol, fileInfo.Size(), contentType, dateHeader)

	if keepAlive {
		headers += "Connection: keep-alive\r\n"
	} else {
		headers += "Connection: close\r\n"
	}
	headers += "\r\n"

	_, err = conn.Write([]byte(headers))
	if err != nil {
		return fmt.Errorf("Error writing headers: %v", err)
	}

	buffer := make([]byte, 8192)
	for {
		n, readErr := file.Read(buffer)
		if n > 0 {
			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				if errors.Is(writeErr, syscall.EPIPE) {
					fmt.Println("Client closed the connection.")
					return nil
				}
				return fmt.Errorf("Error writing file content: %v", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("Error reading file content: %v", readErr)
		}
	}

	return nil
}

// Create and send error response.
func sendError(conn net.Conn, status string, protocol string) {
	body := fmt.Sprintf("<html><body><h1>%s</h1></body></html>", status)
	dateHeader := time.Now().UTC().Format(time.RFC1123)
	headers := fmt.Sprintf("%s %s\r\nContent-Length: %d\r\nDate: %s\r\nConnection: close\r\n\r\n",
		protocol, status, len(body), dateHeader)
	conn.Write([]byte(headers + body))
}

// Uses mime package to get content type.
func getContentType(ext string) string {
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func isKeepAlive(headers map[string]string, protocol string) bool {
	if val, ok := headers["Connection"]; ok {
		return strings.ToLower(val) == "keep-alive"
	}
	return protocol == HTTP11
}
