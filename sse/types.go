package sse

import "encoding/json"

type SSEMessage struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}

func (m *SSEMessage) Unmarshal(payload string) error {
	return json.Unmarshal([]byte(payload), m)
}
