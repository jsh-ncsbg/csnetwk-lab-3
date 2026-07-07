package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/payload"
)

const (
	ServerName   = "Francis"
	ServerNumber = 67
)

func handleClient(shutdown context.CancelFunc, doneChan chan bool, conn net.Conn) {
	svrPayload := payload.NewPayload(ServerName, ServerNumber)

	// read and parse client payload
	cltPayload := &payload.Payload{}
	if err := json.Read(conn, cltPayload); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Println("parsing error: ", err)
		}
		return
	}

	// log
	sum := svrPayload.Number + cltPayload.Number
	log.Printf("received %d from %q", cltPayload.Number, cltPayload.Name)
	log.Printf("%d + %d = %d", svrPayload.Number, cltPayload.Number, sum)

	// reply
	if err := json.Write(conn, svrPayload); err != nil {
		log.Printf("failed to send to client")
	}

	if cltPayload.Number <= 0 || 100 < cltPayload.Number {
		log.Println("client number not in [1, 100]. shutting down..")
		shutdown()
	}

	doneChan <- true
}

func main() {
	serverHost := flag.String("hostname", "0.0.0.0", "server host name")
	serverPort := flag.String("port", "8080", "server port number")

	silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// signal aware context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	log.Println("listening on port", *serverPort)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		<-ctx.Done()
		server.Close()
	})

	// start accepting
	handleCount := 0
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
		handleCount++
		log.Printf("assigning worker %d to %s\n", handleCount, conn.RemoteAddr())

		wg.Go(func() {
			workerId := handleCount
			doneChan := make(chan bool, 1)
			go handleClient(stop, doneChan, conn)

			select {
			case <-ctx.Done():
			case <-doneChan:
			}

			log.Println("closing socket assigned to worker", workerId)
			conn.Close()
		})
	}

	wg.Wait()
	log.Println("server shutdown successfully.")
}
