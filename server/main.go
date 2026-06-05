package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8000"
	NETWORK     = "tcp"

	SERVER_NUMBER = 67
)

func handleClient(connection net.Conn) {
	defer connection.Close()

	decoder := gob.NewDecoder(connection)

	var clientPayload Payload

	if err := decoder.Decode(&clientPayload); err != nil {
		fmt.Println("[Server] Decode error: ", err)
		return
	}

	fmt.Printf("[Server] Received %v\n", clientPayload)
	fmt.Printf("[Server] %d + %d = %d\n", SERVER_NUMBER, clientPayload, SERVER_NUMBER+clientPayload)

	// reply
	if err := sendPayload(connection, SERVER_NUMBER); err != nil {
		fmt.Printf("[Server] Failed to send to client")
	}
}

func main() {
	fmt.Println("Starting...")
	server, err := net.Listen(NETWORK, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Failed to start server: ", err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("[Server] Listening on port", SERVER_PORT)

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("[Server] Connection error. Skipping..")
			continue
		}
		fmt.Println("[Server]", connection.RemoteAddr(), "connected")
		handleClient(connection)
	}
}
