package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

func getClientPayload() (clientPayload protocol.Payload, err error) {
	fmt.Print("[Client] Enter a number in [1, 100]: ")
	_, err = fmt.Scan(&clientPayload)

	return clientPayload, err
}

func getServerPayload(server *net.Conn) (serverPayload protocol.Payload, err error) {
	decoder := gob.NewDecoder(*server)

	if err := decoder.Decode(&serverPayload); err != nil {
		return protocol.Payload(-1), fmt.Errorf("Decode error: %v", err)
	}

	return serverPayload, nil
}

func main() {
	serverHost := flag.String("hostname", "localhost", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")
	network := flag.String("network", "tcp", "network")
	flag.Parse()

	fmt.Println("[Client] Connecting to server...")
	conn, err := net.Dial(*network, *serverHost+":"+strconv.Itoa(*serverPort))
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

	err = protocol.SendPayload(conn, clientPayload)
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
