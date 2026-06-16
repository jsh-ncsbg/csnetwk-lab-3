package main

import (
	"encoding/gob"
	"flag"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

const (
	ServerNumber = 67
)

func handleClient(connection net.Conn) {
	defer connection.Close()

	decoder := gob.NewDecoder(connection)

	var clientPayload protocol.Payload

	if err := decoder.Decode(&clientPayload); err != nil {
		log.Println("Decode error: ", err)
		return
	}

	log.Printf("Received %v", clientPayload)
	log.Printf("%d + %d = %d", ServerNumber, clientPayload, ServerNumber+clientPayload)

	// reply
	if err := protocol.SendPayload(connection, ServerNumber); err != nil {
		log.Printf("Failed to send to client")
	}
}

func main() {
	serverHost := flag.String("hostname", "0.0.0.0", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")
	network := flag.String("network", "tcp", "network")

	silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// logger setup
	log.SetPrefix("[SERVER] ")
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	if *silent {
		log.SetOutput(io.Discard)
	}

	server, err := net.Listen(*network, *serverHost+":"+strconv.Itoa(*serverPort))
	if err != nil {
		log.Fatalln("Failed to start server: ", err)
	}
	defer server.Close()

	log.Println("Listening on port", *serverPort)

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println("Connection error. Skipping..")
			continue
		}
		log.Println(connection.RemoteAddr(), "connected")
		go handleClient(connection)
	}
}
