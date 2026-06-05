# csnetwk-lab-3

### Laboratory 3: Introduction to Sockets Programming
This CSNETWK Laboratory Activity demonstrates basic TCP socket communication between a Client and a Server. The client sends a name and a number to the server, which responds with its own details, allowing both sides to verify the sum of the numbers before safely closing the connection. Sending an out-of-range number gracefully shuts down the server.
## Getting Started

### Client

```bash
cd server
go run server.go
```

### Server

```bash
cd client
go run client.go
```
Alternatively, build and run the executable.