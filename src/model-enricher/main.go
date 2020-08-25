package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/concurrentai/concurrentai-core/src/shared/messaging"
)

// main : Runs the model enricher background service
func main() {
	config := LoadConfig()

	client, closeClient := createPulsarClient(config)
	defer closeClient()

	consumer, closeConsumer := createPulsarConsumer(client, config)
	defer closeConsumer()

	producer, closeProducer := createPulsarProducer(client, config)
	defer closeProducer()

	for {
		if err := HandleNextMessage(consumer, producer); err != nil {
			log.Println(err)
		}
	}
}

// createPulsarClient : Create a Pulsar client
func createPulsarClient(config *Config) (messaging.Client, func()) {
	client, err := messaging.NewPulsarClient(config.PulsarURL)
	if err != nil {
		log.Fatal(err)
	}
	return client, func() {
		client.Close()
	}
}

// createPulsarConsumer : Create a Pulsar consumer
func createPulsarConsumer(client messaging.Client, config *Config) (messaging.Consumer, func()) {
	topic := config.TopicName("model-request")
	subscription := config.SubscriptionName("model-request")
	consumer, err := client.CreateConsumer(topic, subscription)
	if err != nil {
		log.Fatal(err)
	}
	return consumer, func() {
		consumer.Close()
	}
}

// createPulsarProducer : Create a Pulsar producer
func createPulsarProducer(client messaging.Client, config *Config) (messaging.Producer, func()) {
	producer, err := client.CreateProducer(config.TopicName("model-input"))
	if err != nil {
		log.Fatal(err)
	}
	return producer, func() {
		producer.Close()
	}
}

// HandleNextMessage : Enrich request data for a rendezvous request
func HandleNextMessage(consumer messaging.Consumer, producer messaging.Producer) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
	}()

	payload, err := consumer.Receive()
	if err != nil {
		return errors.Wrap(err, "failed to read rendezvous message from consumer")
	}

	if err := producer.Send(payload); err != nil {
		return errors.Wrap(err, "failed to send rendezvous message")
	}

	log.Println("published message: " + string(payload))
	return nil
}
