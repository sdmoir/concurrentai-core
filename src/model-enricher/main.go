package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/concurrent-ai/rendezvous/src/shared/messaging"
)

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

func main() {
	config := LoadConfig()

	client, err := messaging.NewPulsarClient(config.PulsarURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	consumer, err := client.CreateConsumer(config.TopicName("model-request"))
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	producer, err := client.CreateProducer(config.TopicName("model-input"))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	for {
		if err := HandleNextMessage(consumer, producer); err != nil {
			log.Println(err)
		}
	}
}
