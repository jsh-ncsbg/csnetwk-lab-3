package server

import (
	"context"
	"log"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/json"
)

const (
	ServerName   = "Josh"
	ServerNumber = 67
)

func HandleClient(ctxCancelFunc context.CancelFunc, conn net.Conn) {
	defer conn.Close()

	svrPayload := protocol.NewPayload(ServerName, ServerNumber)

	cltPayload := &protocol.Payload{}
	err := json.Read(conn, cltPayload)
	if err != nil {
		log.Println("parsing error: ", err)
		return
	}

	log.Printf("received %d from %v", cltPayload.Number, cltPayload.Name)
	if cltPayload.Number <= 0 || cltPayload.Number > 100 {
		log.Println("client has requested to stop the server.")
		ctxCancelFunc()
		return
	}

	log.Printf("%d + %d = %d", svrPayload.Number, cltPayload.Number, svrPayload.Number+cltPayload.Number)

	// reply
	if err := json.Write(conn, svrPayload); err != nil {
		log.Printf("Failed to send to client")
	}
}
