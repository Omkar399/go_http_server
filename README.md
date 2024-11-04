# go_http_server
A simple golang http web server , to serve concurrent tcp connections.

1. It starts listening to incoming connections on the port(sent from the terminal) as it's now bound to the server.
2. Runs an infinite for loop to accept all incoming connections.
3. Using go routines for func handleConnection(), spawns a new thread for every connection.
4. Function handleConnection() for every new connection if it's HTTP/1.1 it runs an infinite for loop to keep the
    connection open until the dynamic timeout or else close connection after serving.
5. Dynamic timeout- It just uses the number of connections to determine the timeouts, basically the higher the
    number of connections the lower the timeout.
6. It extends the connection timeout every time a new request comes in within the previous timeout to keep the
    persistent tcp connection open.
7. We call func processRequest() inside the handleConnection to check if the request is properly formed, then pass
    it to func serveFile() to actually transmit the file.
8. Serve file checks the file permissions, if it exists and send errors with the appropriate headers (also adds
    keep-alive if its http/1.1), if the file exists we send the files in chunks of 8kb using the buffer.
