package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/payload"
)

func getClientPayload(p *payload.Payload) error {
	fmt.Printf("What is your name? ")
	if _, err := fmt.Scan(&p.Name); err != nil {
		return err
	}

	fmt.Printf("Enter a number in [1, 100]: ")
	if _, err := fmt.Scan(&p.Number); err != nil {
		return err
	}

	return nil
}

func main() {
	serverHost := flag.String("hostname", "localhost", "server host name")
	serverPort := flag.String("port", "8080", "server port number")

	flag.Parse()

	// logger setup
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(log.Lmsgprefix | log.Ltime)

	log.Println("connecting to server...")
	conn, err := net.Dial("tcp", net.JoinHostPort(*serverHost, *serverPort))
	if err != nil {
		log.Fatalln("Error connecting to server: ", err)
	}
	defer conn.Close()

	log.Println("connected to server")

	// client payload
	cltPayload := &payload.Payload{}
	if err := getClientPayload(cltPayload); err != nil {
		log.Fatalln("Input error")
	}

	// send to server
	if err := json.Write(conn, cltPayload); err != nil {
		log.Println("failed to send to server.")
		return
	}

	// receive server payload
	svrPayload := &payload.Payload{}
	if err := json.Read(conn, svrPayload); err != nil {
		log.Println("error communicating with server: ", err)
		return
	}

	log.Printf("received %d from server %q", svrPayload.Number, svrPayload.Name)

	sum := svrPayload.Number + cltPayload.Number
	log.Printf("%d + %d = %d\n", svrPayload.Number, cltPayload.Number, sum)
}
