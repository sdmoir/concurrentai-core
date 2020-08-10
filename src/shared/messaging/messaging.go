package messaging

// Client : A client that can create producers and consumers
type Client interface {
	CreateProducer(topic string) Producer
	CreateConsumer(topic string) Consumer
	Close()
}

// Producer : A message producer that can send messages to a specified topic
type Producer interface {
	Send(payload []byte) error
	Close()
}

// Consumer : A message consumer that can receive messages from a specified topic
type Consumer interface {
	Receive() ([]byte, error)
	Close()
}
