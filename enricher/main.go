package main

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
)

func createPulsarClient(config Config) pulsar.Client {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: config.pulsarURL,
	})

	if err != nil {
		log.Fatal("failed to create pulsar client", err)
	}

	return client
}

func createPulsarConsumer(client pulsar.Client, topic string) pulsar.Consumer {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: fmt.Sprintf("%s-subscription", topic),
		Type:             pulsar.Shared,
	})

	if err != nil {
		log.Fatal("error subscribing to topic:", err)
	}

	return consumer
}

func createPulsarProducer(client pulsar.Client, topic string) pulsar.Producer {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})

	if err != nil {
		log.Fatal("failed to create pulsar producer", err)
	}

	return producer
}

func handleNextMessage(consumer pulsar.Consumer, producer pulsar.Producer) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
	}()

	message, err := consumer.Receive(context.Background())

	if err == nil {
		log.Printf("consumed from topic " + message.Topic() + ": " + string(message.Payload()))

		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: message.Payload(),
		})

		log.Print("published to topic " + producer.Topic())
	}
}

func main() {
	config := LoadConfig()

	client := createPulsarClient(config)
	defer client.Close()

	consumer := createPulsarConsumer(client, config.TopicName("model-request"))
	defer consumer.Close()

	producer := createPulsarProducer(client, config.TopicName("model-input"))
	defer producer.Close()

	for {
		handleNextMessage(consumer, producer)
	}
}
