package server

import (
	"context"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
)

const (
	ServerName   = "Josh"
	ServerNumber = 67
)

func HandleClient(ctxCancelFunc context.CancelFunc, logger logger.Logger, conn net.Conn) {
	defer conn.Close()

	svrPayload := protocol.NewPayload(ServerName, ServerNumber)

	cltPayload := &protocol.Payload{}
	err := json.Read(conn, cltPayload)
	if err != nil {
		logger.Println("parsing error: ", err)
		return
	}

	logger.Printf("received %d from %v\n", cltPayload.Number, cltPayload.Name)
	if cltPayload.Number <= 0 || cltPayload.Number > 100 {
		logger.Println("client has requested to stop the server.")
		ctxCancelFunc()
		return
	}

	logger.Printf("%d + %d = %d\n", svrPayload.Number, cltPayload.Number, svrPayload.Number+cltPayload.Number)

	// reply
	if err := json.Write(conn, svrPayload); err != nil {
		logger.Println("Failed to send to client")
	}
}
