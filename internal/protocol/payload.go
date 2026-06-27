package protocol

type Payload struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func NewPayload(name string, number int) *Payload {
	return &Payload{
		Name:   name,
		Number: number,
	}
}
