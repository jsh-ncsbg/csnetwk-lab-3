# csnetwk-lab-3

### Laboratory 3: Introduction to Sockets Programming

This CSNETWK Laboratory Activity demonstrates basic TCP socket communication between a Client and a Server. The client sends a name and a number to the server, which responds with its own details, allowing both sides to verify the sum of the numbers before safely closing the connection. Sending an out-of-range number gracefully shuts down the server.

## Getting Started

You can run the project using the automated Makefile shortcuts or by executing the Go commands manually.

### Using the Makefile (Recommended)

To run the Server:

```bash
make run-server
```

To run the Client:

```bash
make run-client
```

To build both binaries without running them:

```bash
make build
```

_Note: Compiled binaries will be saved inside the local `bin/` directory._

### Using Manual Go Commands

If you prefer to run the source files directly without using the Makefile:

Server:

```bash
go run ./cmd/server
```

Client:

```bash
go run ./cmd/client
```
