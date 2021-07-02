# Pub/Sub system using Go

A simple Publisher Subscriber system implementation in _Golang_

## Broker APIs

```go
Subscribe(subscriber) // Binds the broker to the subscriber for future communication
```

```go
Unsubscribe(subscriber) // Unbinds the broker to the subscriber
```

```go
AddSubscription(subscriber, topic) // Binds the subscriber to the topic for receiving push messages
```

```go
DeleteSubscription(subscriber, topic) // Unbinds the subscriber to the topic
```

```go
Publish(message, topic) // Pushes the message to all the subscriptions for the topic
```
## Subscriber APIs
```go
CreateTopic(topic) // Adds the topic to the subscriber who calls it
```

```go
DeleteTopic(topic) // Deletes the topic for the subscriber who calls it
```

```go
Ack(subscriber, message) // Initmates the broker that the message has been received and processed. If ACK isn't received in a timeout, broker will retry sending the message
```