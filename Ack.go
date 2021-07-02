package pubsub

//Struct for Ack
type Ack struct {
	payload interface{}
	id      string
}

// GetPayload returns the current message contents
func (a *Ack) GetPayload() interface{} {
	return a.payload
}

// GetID return the subscriber ID of the current ack
func (a *Ack) GetID() string {
	return a.id
}
