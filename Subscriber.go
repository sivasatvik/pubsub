package pubsub

import (
	"encoding/hex"
	"math/rand"
	"sync"
)

// Struct for Subscriber
type Subscriber struct {
	id        string
	destroyed bool
	topics    map[string]bool
	messages  chan *Message
	lock      sync.RWMutex
	ack       chan *Ack
}

// NewSubscriber is a constructor for the Subscriber struct
// which returns a new subscriber
func NewSubscriber() (*Subscriber, error) {
	id := make([]byte, 10)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}

	return &Subscriber{
		id:        hex.EncodeToString(id),
		destroyed: false,
		lock:      sync.RWMutex{},
		topics:    map[string]bool{},
		messages:  make(chan *Message),
		ack:       make(chan *Ack),
	}, nil
}

// CreateTopic adds a new topic to the subscriber
func (s *Subscriber) CreateTopic(topic string) {
	s.lock.Lock()
	s.topics[topic] = true
	s.lock.Unlock()
}

// DeleteTopic removes a topic from the subsciber
func (s *Subscriber) DeleteTopic(topic string) {
	s.lock.Lock()
	delete(s.topics, topic)
	s.lock.Unlock()
}

// Signal sends the message to the subscriber
func (s *Subscriber) Signal(m *Message) *Subscriber {
	s.lock.RLock()
	if !s.destroyed {
		s.messages <- m
	}
	s.lock.RUnlock()
	return s
}

// Ack sends ack signal for the received message
func (s *Subscriber) Ack(m *Message) *Subscriber {
	s.lock.RLock()
	if !s.destroyed {
		a := &Ack{
			payload: m.GetPayload(),
			id:      s.id,
		}
		s.ack <- a
	}
	s.lock.RUnlock()
	return s
}

// GetID returns the ID of the subscriber
func (s *Subscriber) GetID() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.id
}

// GetTopics returns the topics that the
// subscriber has added
func (s *Subscriber) GetTopics() []string {
	s.lock.RLock()
	subscriberTopics := s.topics
	s.lock.RUnlock()

	t := []string{}
	for topic := range subscriberTopics {
		t = append(t, topic)
	}
	return t
}

// GetMessages returns a channel for Message to listen on
func (s *Subscriber) GetMessages() <-chan *Message {
	return s.messages
}

// GetAck returns a channel for Ack to listen on
func (s *Subscriber) GetAck() <-chan *Ack {
	return s.ack
}

// destroy removes the current subscriber
// along with the underlying resources allocated
func (s *Subscriber) destroy() {
	s.lock.Lock()
	s.destroyed = true
	s.lock.Unlock()

	close(s.messages)
}
