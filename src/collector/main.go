package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/concurrent-ai/rendezvous/src/shared/domain"
	"github.com/concurrent-ai/rendezvous/src/shared/messaging"
	"github.com/concurrent-ai/rendezvous/src/shared/sockets"
)

// HandleNextMessage : Receive a rendezvous message and write the model response to the expected socket
func HandleNextMessage(consumer messaging.Consumer, socketWriter sockets.Writer) error {
	payload, err := consumer.Receive()
	if err != nil {
		return errors.Wrap(err, "failed to read rendezvous message from consumer")
	}

	var message *domain.RendezvousMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return errors.Wrap(err, "failed to parse rendezvous message")
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

	consumer, err := client.CreateConsumer(config.TopicName("model-response"))
	if err != nil {
		log.Fatal(err)
	}

	socketWriter := sockets.NewUnixWriter()

	for {
		if err := HandleNextMessage(consumer, socketWriter); err != nil {
			log.Println(err)
		}
	}
}
