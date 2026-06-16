package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

func getClientPayload() (clientPayload protocol.Payload, err error) {
	log.Printf("Enter a number in [1, 100]: ")
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

	// logger setup
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(log.Lmsgprefix | log.Ltime)

	log.Println("Connecting to server...")
	conn, err := net.Dial(*network, *serverHost+":"+strconv.Itoa(*serverPort))

	if err != nil {
		log.Fatalln("Error connecting to server: ", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	clientPayload, err := getClientPayload()
	if err != nil {
		log.Fatalln("Input error")
	}

	err = protocol.SendPayload(conn, clientPayload)
	if err != nil {
		log.Println("Failed to send to server.")
		os.Exit(0)
	}

	if 1 <= clientPayload && clientPayload <= 100 {
		serverPayload, err := getServerPayload(&conn)
		if err != nil {
			log.Println("Error communicating with server: ", err)
		} else {
			sum := serverPayload + clientPayload
			log.Printf("%d + %d = %d\n", serverPayload, clientPayload, sum)
		}
	}
}
