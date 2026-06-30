package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/payload"
)

const (
	ServerName   = "Josh"
	ServerNumber = 67
)

func handleClient(shutdown context.CancelFunc, conn net.Conn) {
	defer conn.Close()

	svrPayload := payload.NewPayload(ServerName, ServerNumber)

	cltPayload := &payload.Payload{}
	err := json.Read(conn, cltPayload)
	if err != nil {
		log.Println("parsing error: ", err)
		return
	}

	log.Printf("received %d from %v", cltPayload.Number, cltPayload.Name)
	log.Printf("%d + %d = %d", svrPayload.Number, cltPayload.Number, svrPayload.Number+cltPayload.Number)

	// reply
	if err := json.Write(conn, svrPayload); err != nil {
		log.Printf("failed to send to client")
	}

	if cltPayload.Number <= 0 || 100 < cltPayload.Number {
		shutdown()
	}
}

func main() {
	serverHost := flag.String("hostname", "0.0.0.0", "server host name")
	serverPort := flag.String("port", "8080", "server port number")

	silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// logger setup
	log.SetPrefix("[SERVER] ")
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	if *silent {
		log.SetOutput(io.Discard)
	}

	server, err := net.Listen("tcp", net.JoinHostPort(*serverHost, *serverPort))
	if err != nil {
		log.Fatalln("failed to start server: ", err)
	}
	defer server.Close()

	log.Println("listening on port", *serverPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	// start accepting
	for {
		conn, err := server.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("stopping accept loop (server shutdown)")
				break
			}

			log.Println("connection error. skipping..")
			continue
		}
		log.Println(conn.RemoteAddr(), "connected")
		go handleClient(cancel, conn)
	}
}
