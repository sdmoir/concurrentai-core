package main

import (
	"log"

	"github.com/concurrent-ai/rendezvous/src/shared/messaging"
)

// HandleNextMessage : Enrich request data for a rendezvous request
func HandleNextMessage(consumer messaging.Consumer, producer messaging.Producer) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
	}()

	payload, err := consumer.Receive()
	if err != nil {
		log.Println(err)
		return
	}

	if err := producer.Send(payload); err != nil {
		log.Println(err)
		return
	}

	log.Println("published message: " + string(payload))
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
		HandleNextMessage(consumer, producer)
	}
}
