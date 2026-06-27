package protocol

import (
	"encoding/json"
	"io"
)

type Payload struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func EncodePayload(w io.Writer, payload Payload) (err error) {
	return json.NewEncoder(w).Encode(payload)
}

func DecodePayload(r io.Reader) (payload Payload, err error) {
	err = json.NewDecoder(r).Decode(&payload)
	return
}
