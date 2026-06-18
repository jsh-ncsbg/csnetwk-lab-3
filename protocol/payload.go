package protocol

import (
	"encoding/gob"
	"io"
)

type Payload struct {
	Name   string
	Number int
}

func SendPayload(w io.Writer, payload Payload) (err error) {
	encoder := gob.NewEncoder(w)
	err = encoder.Encode(payload)
	return
}

func ParsePayload(r io.Reader) (payload Payload, err error) {
	dec := gob.NewDecoder(r)
	err = dec.Decode(&payload)
	return
}
