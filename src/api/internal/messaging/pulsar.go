package messaging

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// PulsarProducer : A producer for sending messages to an Apache Pulsar topic
type PulsarProducer struct {
	pulsarClient   pulsar.Client
	pulsarProducer pulsar.Producer
}

// NewPulsarProducer : Create a new PulsarProducer
func NewPulsarProducer(pulsarURL string, topic string) (*PulsarProducer, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: pulsarURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pulsar client")
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{Topic: topic})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pulsar producer")
	}

	return &PulsarProducer{pulsarClient: client, pulsarProducer: producer}, nil
}

// Send : Send a message payload through the Pulsar producer
func (producer *PulsarProducer) Send(payload []byte) error {
	_, err := producer.pulsarProducer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})

	return err
}

// Close : Close the Pulsar producer and client
func (producer *PulsarProducer) Close() {
	producer.pulsarProducer.Close()
	producer.pulsarClient.Close()
}
