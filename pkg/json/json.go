package json

import (
	"encoding/json"
	"io"
)

func Write(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func Read(r io.Reader, data any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	return decoder.Decode(&data)
}
