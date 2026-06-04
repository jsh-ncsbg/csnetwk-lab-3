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
	SERVER_TYPE = "tcp"
)

type Payload int

func sendPayload(connection net.Conn, payload Payload) (err error) {
	encoder := gob.NewEncoder(connection)
	err = encoder.Encode(payload)
	return
}

func getClientPayload() (clientPayload Payload, err error) {
	fmt.Print("[Client] Enter a number in [1, 100]: ")
	_, err = fmt.Scan(&clientPayload)

	return clientPayload, err
}

func main() {
	fmt.Println("[Client] Connecting to server...")
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("[Client] Error connecting to server: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[Client] Connected to server")

	clientPayload, err := getClientPayload()
	if err != nil {
		fmt.Println("[Client] Input error")
		os.Exit(1)
	}

	err = sendPayload(conn, clientPayload)
	if err != nil {
		fmt.Println("[Client] Failed to send to server.")
		os.Exit(1)
	}

	decoder := gob.NewDecoder(conn)

	var serverPayload Payload
	err = decoder.Decode(&serverPayload)
	if err != nil {
		fmt.Println("[Client] Decode error: ", err)
		os.Exit(1)
	}

	fmt.Printf("[Client] Received %v\n", serverPayload)

	sum := serverPayload + clientPayload
	fmt.Printf("[Client] %d + %d = %d\n", serverPayload, clientPayload, sum)
}
