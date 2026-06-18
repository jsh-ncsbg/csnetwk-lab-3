package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

const (
	ServerName   = "Josh"
	ServerNumber = 67
)

func handleClient(connection net.Conn) {
	defer connection.Close()

	serverPayload := protocol.Payload{
		Name:   ServerName,
		Number: ServerNumber,
	}

	clientPayload, err := protocol.ParsePayload(connection)
	if err != nil {
		log.Println("parsing error: ", err)
		return
	}

	log.Printf("received %d from %v", clientPayload.Number, clientPayload.Name)
	log.Printf("%d + %d = %d", serverPayload.Number, clientPayload.Number, serverPayload.Number+clientPayload.Number)

	// reply
	if err := protocol.SendPayload(connection, serverPayload); err != nil {
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
