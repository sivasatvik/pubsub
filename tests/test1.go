package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"pubsub"
	"time"
)

// Test for multiple subscribers on the same topic
func main() {
	broker := pubsub.NewBroker()

	subscriber1, err := broker.Subscribe()
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	subscriber2, err := broker.Subscribe()
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	broker.AddSubscription(subscriber1, "ABCDE")
	broker.AddSubscription(subscriber2, "ABCDE")

	fmt.Println("Subscribers for ABCDE: ", broker.Subscribers("ABCDE"))

	fmt.Println("SubscriptionID: ", subscriber1.GetID(), ", topics subscribed: ", subscriber1.GetTopics())

	ch1 := subscriber1.GetMessages()
	ch2 := subscriber2.GetMessages()

	go send(broker)
	go receive(subscriber1.GetID(), ch1)
	go receive(subscriber2.GetID(), ch2)

	fmt.Scanln()
	fmt.Println("done")
}

func createMessage() string {
	m := make([]byte, 20)
	if _, err := rand.Read(m); err != nil {
		return ""
	}
	return hex.EncodeToString(m)
}

func send(broker *pubsub.Broker) {
	fmt.Println("----Sending----")
	for {
		m := createMessage()
		fmt.Printf("On topic: ABCDE, sending the message: %v\n", m)
		broker.Publish(m, "ABCDE")
		time.Sleep(time.Second)
	}
}

func receive(id string, ch <-chan *pubsub.Message) {
	fmt.Printf("----Subscriber %s, receiving----\n", id)
	for {
		if msg, ok := <-ch; ok {
			fmt.Printf("Subscriber %s, on topic: %s, received the message: %v\n", id, msg.GetTopic(), msg.GetPayload())
		}
	}
}
