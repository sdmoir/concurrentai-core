package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/concurrentai/concurrentai-core/src/shared/domain"
	"github.com/concurrentai/concurrentai-core/src/shared/messaging"
	"github.com/concurrentai/concurrentai-core/src/shared/sockets"
)

// main : Runs the rendezvous collector background service
func main() {
	config := LoadConfig()

	client, closeClient := createPulsarClient(config)
	defer closeClient()

	consumer, closeConsumer := createPulsarConsumer(client, config)
	defer closeConsumer()

	socketWriter := sockets.NewUnixWriter()

	for {
		if err := HandleNextMessage(consumer, socketWriter, config); err != nil {
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
	topic := config.TopicName("model-response")
	subscription := config.SubscriptionName("model-response")
	consumer, err := client.CreateConsumer(topic, subscription)
	if err != nil {
		log.Fatal(err)
	}
	return consumer, func() {
		consumer.Close()
	}
}

// HandleNextMessage : Receive a rendezvous message and write the model response to the expected socket
func HandleNextMessage(consumer messaging.Consumer, socketWriter sockets.Writer, config *Config) error {
	payload, err := consumer.Receive()
	if err != nil {
		return errors.Wrap(err, "failed to read rendezvous message from consumer")
	}

	var message *domain.RendezvousMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return errors.Wrap(err, "failed to parse rendezvous message")
	}

	if message.ResponseModelID != config.ActiveModelID {
		return nil
	}

	socketAddress := fmt.Sprintf("/sockets/%s.sock", message.ID)
	log.Println("Writing response to " + socketAddress)
	if err := socketWriter.Write(socketAddress, []byte(message.ResponseData)); err != nil {
		return errors.Wrap(err, "failed to write rendezvous message response data to socket")
	}

	return nil
}
