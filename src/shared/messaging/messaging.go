package messaging

// Producer : A message producer that can send messages to a specified topic
type Producer interface {
	Send(payload []byte) error
	Close()
}
