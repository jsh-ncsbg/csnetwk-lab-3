package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	NETWORK     = "tcp"
)

func getClientPayload() (clientPayload Payload, err error) {
	fmt.Print("[Client] Enter a number in [1, 100]: ")
	_, err = fmt.Scan(&clientPayload)

	return clientPayload, err
}

func getServerPayload(server *net.Conn) (serverPayload Payload, err error) {
	decoder := gob.NewDecoder(*server)

	if err := decoder.Decode(&serverPayload); err != nil {
		return Payload(-1), fmt.Errorf("Decode error: %v", err)
	}

	return serverPayload, nil
}

func main() {
	fmt.Println("[Client] Connecting to server...")
	conn, err := net.Dial(NETWORK, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("[Client] Error connecting to server: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[Client] Connected to server")

	clientPayload, err := getClientPayload()
	if err != nil {
		fmt.Println("[Client] Input error")
		return
	}

	err = sendPayload(conn, clientPayload)
	if err != nil {
		fmt.Println("[Client] Failed to send to server.")
		return
	}

	if 1 <= clientPayload && clientPayload <= 100 {
		serverPayload, err := getServerPayload(&conn)
		if err != nil {
			fmt.Println("[Client] Error communicating with server")
		} else {
			sum := serverPayload + clientPayload
			fmt.Printf("[Client] %d + %d = %d\n", serverPayload, clientPayload, sum)
		}
	}
}
