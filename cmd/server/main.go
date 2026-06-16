package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

const (
	ServerHost = "localhost"
	ServerPort = 8000
	NETWORK    = "tcp"

	ServerNumber = 67
)

func handleClient(connection net.Conn) {
	defer connection.Close()

	decoder := gob.NewDecoder(connection)

	var clientPayload protocol.Payload

	if err := decoder.Decode(&clientPayload); err != nil {
		fmt.Println("[Server] Decode error: ", err)
		return
	}

	fmt.Printf("[Server] Received %v\n", clientPayload)
	fmt.Printf("[Server] %d + %d = %d\n", ServerNumber, clientPayload, ServerNumber+clientPayload)

	// reply
	if err := protocol.SendPayload(connection, ServerNumber); err != nil {
		fmt.Printf("[Server] Failed to send to client")
	}
}

func main() {
	ServerHost := flag.String("hostname", "0.0.0.0", "server host name")
	ServerPort := flag.Int("port", 8000, "server port number")
	flag.Parse()

	fmt.Println("Starting...")
	server, err := net.Listen(NETWORK, fmt.Sprintf("%s:%d", *ServerHost, *ServerPort))
	if err != nil {
		fmt.Println("Failed to start server: ", err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("[Server] Listening on port", *ServerPort)

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("[Server] Connection error. Skipping..")
			continue
		}
		fmt.Println("[Server]", connection.RemoteAddr(), "connected")
		go handleClient(connection)
	}
}
