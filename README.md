# csnetwk-lab-3

### Laboratory 3: Introduction to Sockets Programming

This CSNETWK Laboratory Activity demonstrates basic TCP socket communication between a Client and a Server. The client sends a name and a number to the server, which responds with its own details, allowing both sides to verify the sum of the numbers before safely closing the connection. Sending an out-of-range number gracefully shuts down the server.

## Getting Started

### Server

```bash
go run ./cmd/server
```

Flags:

| Flag | Default | Description |
| --- | --- | --- |
| `-hostname` | `0.0.0.0` | Host name/address to listen on |
| `-port` | `8080` | Port number to listen on |
| `-silent` | `false` | Disable all log output |

### Client

```bash
go run ./cmd/client
```

Flags:

| Flag | Default | Description |
| --- | --- | --- |
| `-hostname` | `localhost` | Server host name to connect to |
| `-port` | `8080` | Server port number to connect to |

### Using the Makefile

A `Makefile` is included to build and run both programs:

```bash
make build         # build both server and client into ./bin
make build-server   # build only the server
make build-client   # build only the client

make run-server     # build (if needed) and run the server
make run-client      # build (if needed) and run the client
```

Alternatively, build and run the executable directly.

## Simulator

`cmd/simulator` spins up many concurrent fake clients against a running server to benchmark it. It currently connects to `localhost:8080`, so start the server first.

```bash
go run ./cmd/simulator
```

Flags:

| Flag | Default | Description |
| --- | --- | --- |
| `-concurrent` | `50` | Number of concurrent simulated users |
| `-total` | `1000` | Total number of users to simulate |
| `-delay` | `0` | Delay (in milliseconds) before each simulated client replies, to mimic lag |
