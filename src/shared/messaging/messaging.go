package messaging

// Client : A client that can create producers and consumers
type Client interface {
	CreateProducer(topic string) (Producer, error)
	CreateConsumer(topic string, subscriptionName string) (Consumer, error)
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
