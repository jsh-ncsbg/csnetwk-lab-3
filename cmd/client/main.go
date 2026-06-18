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

	// send to server
	err = protocol.SendPayload(conn, clientPayload)
	if err != nil {
		log.Println("Failed to send to server.")
		os.Exit(0)
	}

	if 1 <= clientPayload.Number && clientPayload.Number <= 100 {
		serverPayload, err := protocol.ParsePayload(conn)
		if err != nil {
			log.Println("Error communicating with server: ", err)
		} else {
			log.Printf("received %d from server %v", serverPayload.Number, serverPayload.Name)

			sum := serverPayload.Number + clientPayload.Number
			log.Printf("%d + %d = %d\n", serverPayload.Number, clientPayload.Number, sum)
		}
	}
}
