package main

import (
	"log"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
)

const (
	ServerName   = "Josh"
	ServerNumber = 67
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	svrPayload := protocol.Payload{
		Name:   ServerName,
		Number: ServerNumber,
	}

	var cltPayload protocol.Payload
	err := json.Read(conn, &cltPayload)
	if err != nil {
		log.Println("parsing error: ", err)
		return
	}

	log.Printf("received %d from %v", cltPayload.Number, cltPayload.Name)
	log.Printf("%d + %d = %d", svrPayload.Number, cltPayload.Number, svrPayload.Number+cltPayload.Number)

	// reply
	if err := json.Write(conn, svrPayload); err != nil {
		log.Printf("Failed to send to client")
	}
}
