package sse

import "sync"

type Stream struct {
	subscribers map[string]*Subscriber
	rwm         *sync.RWMutex
}

func NewStream() *Stream {
	return &Stream{
		subscribers: make(map[string]*Subscriber),
		rwm:         &sync.RWMutex{},
	}
}

func (s *Stream) AddSubscriber(sub *Subscriber) {
	s.rwm.Lock()
	defer s.rwm.Unlock()

	s.subscribers[sub.uid] = sub
}

func (s *Stream) RemoveSubscriber(uid string) {
	s.rwm.Lock()
	defer s.rwm.Unlock()

	v, ok := s.subscribers[uid]
	if !ok {
		return
	}

	close(v.dataChan) // close chan
	delete(s.subscribers, uid)
}

func (s *Stream) Broadcast(payload string) {
	s.rwm.RLock()
	defer s.rwm.RUnlock()

	for _, v := range s.subscribers {
		v.dataChan <- payload
	}
}

func (s *Stream) CanBeRemoved() bool {
	return len(s.subscribers) == 0
}
