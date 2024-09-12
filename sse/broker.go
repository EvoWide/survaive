package sse

import (
	"encoding/json"
	"survaive/bus"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Broker struct {
	streams map[string]*Stream
	rwm     *sync.RWMutex

	bus bus.Bus
}

func NewBroker() *Broker {
	return &Broker{
		streams: make(map[string]*Stream),
		rwm:     &sync.RWMutex{},
	}
}

func (b *Broker) AttachRedisBus(client *redis.Client, channel string) {
	b.bus = bus.NewRedisBus(client, channel)
	c := b.bus.Subscribe()

	go func() {
		for {
			payload := <-c

			msg := &SSEMessage{}
			err := msg.Unmarshal(payload)

			// failed to decode payload!
			if err != nil {
				continue
			}

			b.BroadcastLocally(msg.Channel, msg.Payload)
		}
	}()
}

func (b *Broker) AddSubscriber(name string, uid string) chan string {
	sub := NewSubscriber(uid, make(chan string))
	stream := b.GetOrCreateStream(name)

	stream.AddSubscriber(sub)
	return sub.dataChan
}

func (b *Broker) RemoveSubscriber(name string, uid string) {
	stream := b.GetStream(name)
	if stream == nil {
		return
	}

	stream.RemoveSubscriber(uid)

	if stream.CanBeRemoved() {
		b.DeleteStream(name)
	}
}

func (b *Broker) BroadcastLocally(name string, payload string) {
	if stream := b.GetStream(name); stream != nil {
		stream.Broadcast(payload)
	}
}

func (b *Broker) Broadcast(c string, payload string) error {
	if b.bus != nil {
		out, err := json.Marshal(&SSEMessage{c, payload})
		if err != nil {
			return err
		}

		if err = b.bus.Publish(string(out)); err != nil {
			return err
		}
	}

	b.BroadcastLocally(c, payload)
	return nil
}

func (b *Broker) GetOrCreateStream(name string) *Stream {
	stream := b.GetStream(name)
	if stream != nil {
		return stream
	}

	return b.CreateStream(name)
}

func (b *Broker) GetStream(name string) *Stream {
	b.rwm.RLock()
	defer b.rwm.RUnlock()

	return b.streams[name]
}

func (b *Broker) CreateStream(name string) *Stream {
	b.rwm.Lock()
	defer b.rwm.Unlock()

	stream := NewStream()
	b.streams[name] = stream

	return stream
}

func (b *Broker) DeleteStream(name string) {
	b.rwm.Lock()
	defer b.rwm.Unlock()

	delete(b.streams, name)
}
