package protocol

import (
	"encoding/gob"
	"net"
)

type Payload int

func SendPayload(connection net.Conn, payload Payload) (err error) {
	encoder := gob.NewEncoder(connection)
	err = encoder.Encode(payload)
	return
}
