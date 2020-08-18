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

func main() {
	config := LoadConfig()

	client, err := messaging.NewPulsarClient(config.PulsarURL)
	if err != nil {
		log.Fatal(err)
	}

	topic := config.TopicName("model-response")
	subscription := config.SubscriptionName("model-response")
	consumer, err := client.CreateConsumer(topic, subscription)
	if err != nil {
		log.Fatal(err)
	}

	socketWriter := sockets.NewUnixWriter()

	for {
		if err := HandleNextMessage(consumer, socketWriter, config); err != nil {
			log.Println(err)
		}
	}
}
