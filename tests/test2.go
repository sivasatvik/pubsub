package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"pubsub"
	"time"
)

// Test for one subscriber for multiple topics
func main() {
	broker := pubsub.NewBroker()

	subscriber1, err := broker.Subscribe()
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	broker.AddSubscription(subscriber1, "ABCDE")
	broker.AddSubscription(subscriber1, "STUVW")
	broker.AddSubscription(subscriber1, "LMNOP")

	fmt.Println("Subscribers for ABCDE: ", broker.Subscribers("ABCDE"))

	fmt.Println("SubscriptionID: ", subscriber1.GetID(), ", topics subscribed: ", subscriber1.GetTopics())

	ch1 := subscriber1.GetMessages()

	go send("ABCDE", broker)
	go send("STUVW", broker)
	go send("LMNOP", broker)
	go receive(subscriber1, ch1)

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

func send(topic string, broker *pubsub.Broker) {
	fmt.Println("----Sending----")
	for {
		m := createMessage()
		fmt.Printf("On topic: %s, sending the message: %v\n", topic, m)
		broker.Publish(m, topic)
		time.Sleep(2 * time.Second)
		break
	}
}

func receive(sub *pubsub.Subscriber, ch <-chan *pubsub.Message) {
	fmt.Printf("----Subscriber %s, receiving----\n", sub.GetID())
	for {
		if msg, ok := <-ch; ok {
			fmt.Printf("Subscriber %v, on topic: %s, receiving the message: %s\n", sub.GetID(), msg.GetTopic(), msg.GetPayload())
			sub.Ack(msg)
		}
	}
}
