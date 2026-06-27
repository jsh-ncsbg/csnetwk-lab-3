package client

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/protocol"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
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

func Serve(logger logger.Logger, w io.Writer, r io.Reader, rw io.ReadWriter) {
	cltPayload, err := getClientPayload(w, r)
	if err != nil {
		logger.Println("Input error")
		return
	}

	// send to server
	err = json.Write(rw, cltPayload)
	if err != nil {
		logger.Println("failed to send to server: ", err)
		return
	}

	if 1 <= cltPayload.Number && cltPayload.Number <= 100 {
		svrPayload := &protocol.Payload{}

		err := json.Read(rw, svrPayload)
		if err != nil {
			logger.Println("error communicating with server: ", err)
		} else {
			logger.Printf("received %d from server %v\n", svrPayload.Number, svrPayload.Name)

			sum := svrPayload.Number + cltPayload.Number
			logger.Printf("%d + %d = %d\n", svrPayload.Number, cltPayload.Number, sum)
		}
	}
}
