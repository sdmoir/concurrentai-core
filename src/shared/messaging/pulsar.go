package messaging

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// PulsarClient : A client for creating Pulsar producers and consumers
type PulsarClient struct {
	internalClient pulsar.Client
}

// NewPulsarClient : Create a new Pulsar client
func NewPulsarClient(pulsarURL string) (*PulsarClient, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: pulsarURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pulsar client")
	}
	return &PulsarClient{internalClient: client}, nil
}

// CreateProducer : Create a new Pulsar producer
func (client *PulsarClient) CreateProducer(topic string) (*PulsarProducer, error) {
	producer, err := client.internalClient.CreateProducer(pulsar.ProducerOptions{Topic: topic})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pulsar producer")
	}
	return &PulsarProducer{internalProducer: producer}, nil
}

// CreateConsumer : Create a new Pulsar consumer
func (client *PulsarClient) CreateConsumer(topic string) (*PulsarConsumer, error) {
	consumer, err := client.internalClient.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: fmt.Sprintf("%s-subscription", topic),
		Type:             pulsar.Shared,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pulsar consumer")
	}
	return &PulsarConsumer{internalConsumer: consumer}, nil
}

// Close : Close the Pulsar client
func (client *PulsarClient) Close() {
	client.internalClient.Close()
}

// PulsarProducer : A producer for sending messages to an Apache Pulsar topic
type PulsarProducer struct {
	internalProducer pulsar.Producer
}

// Send : Send a message payload through the Pulsar producer
func (producer *PulsarProducer) Send(payload []byte) error {
	_, err := producer.internalProducer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}

// Close : Close the Pulsar producer
func (producer *PulsarProducer) Close() {
	producer.internalProducer.Close()
}

// PulsarConsumer : A consumer for receiving messages from an Apache Pulsar topic
type PulsarConsumer struct {
	internalConsumer pulsar.Consumer
}

// Receive : Receive a message through the Pulsar consumer
func (consumer *PulsarConsumer) Receive() ([]byte, error) {
	message, err := consumer.internalConsumer.Receive(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to receive message")
	}
	return message.Payload(), nil
}

// Close : Close the Pulsar consumer
func (consumer *PulsarConsumer) Close() {
	consumer.internalConsumer.Close()
}
