package pubsub

// Struct for Message
type Message struct {
	topic   string
	payload interface{}
}

// GetPayload returns the current message contents
func (m *Message) GetPayload() interface{} {
	return m.payload
}

// GetTopic returns the topic in which the message is sent
func (m *Message) GetTopic() string {
	return m.topic
}
