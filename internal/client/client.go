package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/json"
)

func getClientPayload(w io.Writer, r io.Reader) (*protocol.Payload, error) {
	bufReader := bufio.NewReader(r)

	// client name
	fmt.Fprintf(w, "What is your name? ")
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	cltName := strings.TrimSpace(line)

	// client number
	fmt.Fprintf(w, "Enter a number in [1, 100]: ")
	line, err = bufReader.ReadString('\n')
	line = strings.TrimSpace(line)
	cltNumber, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	return protocol.NewPayload(cltName, cltNumber), nil
}

func Serve(w io.Writer, r io.Reader, rw io.ReadWriter) {
	cltPayload, err := getClientPayload(w, r)
	if err != nil {
		log.Fatalln("Input error")
	}

	// send to server
	err = json.Write(rw, cltPayload)
	if err != nil {
		fmt.Fprint(w, "failed to send to server: ", err)
		return
	}

	if 1 <= cltPayload.Number && cltPayload.Number <= 100 {
		svrPayload := &protocol.Payload{}

		err := json.Read(rw, svrPayload)
		if err != nil {
			fmt.Fprintln(w, "error communicating with server: ", err)
		} else {
			fmt.Fprintf(w, "received %d from server %v", svrPayload.Number, svrPayload.Name)

			sum := svrPayload.Number + cltPayload.Number
			fmt.Fprintf(w, "%d + %d = %d\n", svrPayload.Number, cltPayload.Number, sum)
		}
	}
}
