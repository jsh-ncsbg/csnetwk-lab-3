package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

func getClientPayload() (payload protocol.Payload, err error) {
	log.Printf("Enter a number in [1, 100]: ")

	var number int
	_, err = fmt.Scan(&number)

	payload = protocol.Payload(number)

	return
}

func main() {
	serverHost := flag.String("hostname", "localhost", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")
	network := flag.String("network", "tcp", "network")
	flag.Parse()

	// logger setup
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(log.Lmsgprefix | log.Ltime)

	log.Println("Connecting to server...")
	conn, err := net.Dial(*network, *serverHost+":"+strconv.Itoa(*serverPort))

	if err != nil {
		log.Fatalln("Error connecting to server: ", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	clientPayload, err := getClientPayload()
	if err != nil {
		log.Fatalln("Input error")
	}

	err = protocol.SendPayload(conn, clientPayload)
	if err != nil {
		log.Println("Failed to send to server.")
		os.Exit(0)
	}

	if 1 <= clientPayload && clientPayload <= 100 {
		serverPayload, err := protocol.ParsePayload(conn)
		if err != nil {
			log.Println("Error communicating with server: ", err)
		} else {
			sum := serverPayload + clientPayload
			log.Printf("%d + %d = %d\n", serverPayload, clientPayload, sum)
		}
	}
}
