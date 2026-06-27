package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
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

func main() {
	serverHost := flag.String("host", "", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// logger setup
	log.SetPrefix("[SERVER] ")
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	if *silent {
		log.SetOutput(io.Discard)
	}

	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}
	server := server.NewServer(serverConfig)

	err := server.Run()
	if err != nil {
		log.Fatal("failed to start server", err)
	}
	defer server.Close()
	log.Print("listening on ", serverConfig.Host, ":", serverConfig.Port)

	server.HandleFunc(handleClient)
	server.Start()
}
