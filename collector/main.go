package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"net"
	"os"
	"time"
)

var config = LoadConfig()

func handleNextMessage(consumer pulsar.Consumer) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
		// fmt.Println("here")
	}()

	message, err := consumer.Receive(context.Background())

	if err == nil {
		log.Printf("consumed from topic %s at offset %v: "+
			string(message.Payload()), message.Topic())

		writeRendezvousResponse(message)
	}
}

type rendezvousModelResponse struct {
	id            string
	events        map[string]interface{}
	body          map[string]interface{}
	modelID       string
	modelResponse string
}

func writeRendezvousResponse(message pulsar.Message) {
	var modelResponse *rendezvousModelResponse
	if err := json.Unmarshal(message.Payload(), &modelResponse); err != nil {
		log.Fatal(err)
	}

	if modelResponse.modelID != config.activeModelID {
		log.Println("Received shadow model response" + modelResponse.modelResponse)
		return
	}

	log.Println("Received active model response" + modelResponse.modelResponse)

	socketAddress := fmt.Sprintf("/sockets/%s.sock", modelResponse.id)

	for i := 0; i <= 300; i++ {
		if _, err := os.Stat(socketAddress); os.IsNotExist(err) {
			if i < 300 {
				time.Sleep(10 * time.Millisecond)
			} else {
				return
			}
		}
	}

	connection, error := net.Dial("unix", socketAddress)
	if error != nil {
		log.Println("dial error:", error)
	}
	defer connection.Close()

	response := fmt.Sprintf("{ \"results\": %s }", modelResponse.modelResponse)

	_, error = connection.Write([]byte(response))
	if error != nil {
		log.Fatal("write error:", error)
	}
}

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: config.pulsarURL,
	})

	if err != nil {
		log.Fatal("failed to create pulsar client", err)
	}

	defer client.Close()

	topic := config.TopicName("model-response")

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: fmt.Sprintf("%s-subscription", topic),
		Type:             pulsar.Shared,
	})

	if err != nil {
		log.Fatal("error subscribing to topic:", err)
	}

	defer consumer.Close()

	for {
		handleNextMessage(consumer)
	}
}
