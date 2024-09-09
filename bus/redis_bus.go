package bus

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisBus struct {
	BusId   string
	Channel string
	Client  *redis.Client
}

func NewRedisBus(client *redis.Client, channel string) *RedisBus {
	busId := generateId()
	return &RedisBus{busId, channel, client}
}

func (b *RedisBus) Publish(payload string) {
	b.Client.Publish(context.Background(), b.Channel, &BusMessage{b.BusId, payload}) // todo: handle result
}

func (b *RedisBus) Subscribe() chan string {
	pubsub := b.Client.Subscribe(context.Background(), b.Channel)
	defer pubsub.Close()

	c := make(chan string)

	go func() {
		for {
			raw := <-pubsub.Channel()

			msg := &BusMessage{}
			err := json.Unmarshal([]byte(raw.Payload), msg)

			if err != nil {
				continue
			}

			if msg.BusId == b.BusId {
				continue // we don't want to process if bus id is this one
			}

			c <- msg.Payload
		}
	}()

	return c
}
