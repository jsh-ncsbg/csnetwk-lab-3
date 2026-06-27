package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
)

func getClientPayload(r io.Reader) (*protocol.Payload, error) {
	bufReader := bufio.NewReader(r)

	// client name
	fmt.Printf("What is your name? ")
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	cltName := strings.TrimSpace(line)

	// client number
	fmt.Printf("Enter a number in [1, 100]: ")
	line, err = bufReader.ReadString('\n')
	line = strings.TrimSpace(line)
	cltNumber, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	return protocol.NewPayload(cltName, cltNumber), nil
}

func serve(r io.Reader, rw io.ReadWriter) {
	cltPayload, err := getClientPayload(r)
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
