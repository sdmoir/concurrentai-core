package main

import (
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

// CreatePulsarClient : Create a client for interacting with an Apache Pulsar cluster
func CreatePulsarClient(config Config) pulsar.Client {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: config.pulsarURL,
	})

	if err != nil {
		log.Fatal("failed to create pulsar client", err)
	}

	return client
}

// CreatePulsarProducer : Create a producer for sending messages to an Apache Pulsar topic
func CreatePulsarProducer(client pulsar.Client, topic string) pulsar.Producer {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})

	if err != nil {
		log.Fatal("failed to create pulsar producer", err)
	}

	return producer
}
