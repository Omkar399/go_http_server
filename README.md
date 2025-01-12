# GoLang HTTP Server

A **simple and efficient HTTP server** implemented in GoLang to handle **concurrent TCP connections**. This project demonstrates the use of Go's concurrency features to manage multiple client connections while maintaining performance and scalability.

## Features

- **Concurrent Connection Handling**: Uses goroutines for managing multiple client connections simultaneously.
- **Dynamic Timeout Management**: Adjusts connection timeouts dynamically based on the number of active connections.
- **Persistent Connections**: Supports HTTP/1.1 persistent connections with timeout extensions for active requests.
- **Efficient File Serving**: Serves files in 8KB chunks with proper header management, including error handling and permission checks.

## How It Works

1. **Server Initialization**:
   - The server listens on a port specified via the terminal.
   - An infinite loop accepts incoming connections.

2. **Connection Handling**:
   - Each connection spawns a new goroutine using the `handleConnection()` function.
   - For HTTP/1.1 requests, the connection remains open until the dynamic timeout expires or the client disconnects.

3. **Dynamic Timeout**:
   - Timeout duration decreases as the number of active connections increases.
   - Timeout is extended if a new request arrives within the current timeout window.

4. **Request Processing**:
   - Requests are validated using `processRequest()` to ensure proper formatting.
   - Valid requests are passed to `serveFile()` for file transmission.

5. **File Serving**:
   - Checks file permissions and existence before serving.
   - Sends files in 8KB chunks with appropriate headers (e.g., `Keep-Alive` for HTTP/1.1).
   - Handles errors gracefully with descriptive responses.

## How to Run This?

1. Ensure you have Go installed on your machine:
   - If you use Homebrew, simply run:
     ```
     brew install go
     ```

2. Run the server from the terminal:
    ```
    go run http_web_server.go -document_root path/to/documentRoot
    ```
The server will now listen for incoming connections on `localhost:8080`.
