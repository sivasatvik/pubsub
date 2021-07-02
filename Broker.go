package pubsub

import (
	"sync"
	"time"
)

// Subscribers map to their id for easy retrieval
type Subscribers map[string]*Subscriber

// Struct for the broker between publisher and subscriber
// Broker takes care of the subscriptions and publishing
type Broker struct {
	topics      map[string]Subscribers
	topicsLock  sync.RWMutex
	subscribers Subscribers
	subsLock    sync.RWMutex
}

// NewBroker is a constructor for the Broker struct
// which returns a new broker
func NewBroker() *Broker {
	return &Broker{
		topics:      map[string]Subscribers{},
		topicsLock:  sync.RWMutex{},
		subscribers: Subscribers{},
		subsLock:    sync.RWMutex{},
	}
}

// Subscribe binds the subscriber to the broker for
// further services
func (b *Broker) Subscribe() (*Subscriber, error) {
	s, err := NewSubscriber()

	if err != nil {
		return nil, err
	}

	b.subsLock.Lock()
	b.subscribers[s.GetID()] = s
	b.subsLock.Unlock()

	return s, nil
}

// Unsubscribe unbinds the subscriber from the broker
// for opting out of any subscriptions in the future
func (b *Broker) Unsubscribe(s *Subscriber) {
	s.destroy()
	b.subsLock.Lock()
	b.DeleteSubscription(s, s.GetTopics()...)
	delete(b.subscribers, s.id)
	defer b.subsLock.Unlock()
}

// AddSubscription adds subscription from subscriber to topic
func (b *Broker) AddSubscription(s *Subscriber, topics ...string) {
	b.topicsLock.Lock()
	defer b.topicsLock.Unlock()

	for _, topic := range topics {
		if nil == b.topics[topic] {
			b.topics[topic] = Subscribers{}
		}
		s.CreateTopic(topic)
		b.topics[topic][s.id] = s
	}
}

// DeleteSubscription deletes subscription from subscriber to topic
func (b *Broker) DeleteSubscription(s *Subscriber, topics ...string) {
	for _, topic := range topics {
		b.topicsLock.Lock()
		if nil == b.topics[topic] {
			b.topicsLock.Unlock()
			continue
		}
		delete(b.topics[topic], s.id)
		b.topicsLock.Unlock()
		s.DeleteTopic(topic)
	}
}

// Publish pushed the message (payload) to all the subscriptions
// to the topics
func (b *Broker) Publish(payload interface{}, topics ...string) {
	for _, topic := range topics {
		if b.Subscribers(topic) < 1 {
			continue
		}
		b.topicsLock.RLock()
		for _, s := range b.topics[topic] {
			m := &Message{
				topic:   topic,
				payload: payload,
			}
			ch := s.GetAck()
			go (func(s *Subscriber) {
			RETRY:
				s.Signal(m)
				select {
				case msg := <-ch:
					// ACK has been received for the message sent, so move ahead
					//fmt.Printf("Received ack from the subscriber: %s\n", msg.GetID())
					_ = msg
				case <-time.After(5 * time.Second):
					//fmt.Printf("Didn't receive ack from the subscriber: %s, retrying\n", s.GetID())
					// ACK hasn't been received in the time limit. Retry sending the message
					goto RETRY
				}
			})(s)
		}
		b.topicsLock.RUnlock()
	}
}

// Subscribers returns all the current subscriptions
// to the topic
func (b *Broker) Subscribers(topic string) int {
	b.topicsLock.RLock()
	defer b.topicsLock.RUnlock()
	return len(b.topics[topic])
}

// GetTopics returns all the current topics
func (b *Broker) GetTopics() []string {
	b.topicsLock.RLock()
	brokerTopics := b.topics
	b.topicsLock.RUnlock()

	t := []string{}
	for topic := range brokerTopics {
		t = append(t, topic)
	}
	return t
}
