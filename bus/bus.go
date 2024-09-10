package bus

import (
	"encoding/json"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type BusMessage struct {
	BusId   string `json:"busId"`
	Payload string `json:"payload"`
}

func (m *BusMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type Bus interface {
	Publish(message string)
	Subscribe() chan string
}

// move that somewhere else?
func generateId() string {
	busId, err := gonanoid.New()
	if err != nil {
		panic("bus: failed to generateId")
	}
	return busId
}
