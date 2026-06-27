package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/client"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
)

func getClientPayload() (payload protocol.Payload, err error) {
	var number int
	var name string

	fmt.Printf("What is your name? ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return protocol.Payload{}, err
	}

	fmt.Printf("Enter a number in [1, 100]: ")
	_, err = fmt.Scan(&number)
	if err != nil {
		return protocol.Payload{}, err
	}

	return protocol.Payload{Name: name, Number: number}, err
}

func serve(r io.Reader, rw io.ReadWriter) {
	cltPayload, err := getClientPayload()
	if err != nil {
		log.Fatalln("Input error")
	}

	// send to server
	err = json.Write(rw, cltPayload)
	if err != nil {
		log.Println("Failed to send to server.")
		os.Exit(0)
	}

	if 1 <= cltPayload.Number && cltPayload.Number <= 100 {
		var svrPayload protocol.Payload

		err := json.Read(rw, &svrPayload)
		if err != nil {
			log.Println("Error communicating with server: ", err)
		} else {
			log.Printf("received %d from server %v", svrPayload.Number, svrPayload.Name)

			sum := svrPayload.Number + cltPayload.Number
			log.Printf("%d + %d = %d\n", svrPayload.Number, cltPayload.Number, sum)
		}
	}
}

func main() {
	serverHost := flag.String("hostname", "localhost", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	flag.Parse()

	// logger setup
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(log.Lmsgprefix | log.Ltime)

	log.Println("Connecting to server...")
	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}

	client := client.NewClient(os.Stdin, serverConfig)
	client.HandleFunc(serve)

	if err := client.Start(); err != nil {
		log.Fatalln("Error connecting to server: ", err)
	}

	log.Println("Connected to server")
	defer client.Close()
}
