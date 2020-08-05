package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var config = LoadConfig()

func handleNextMessage(consumer pulsar.Consumer) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic occurred: %s", err)
		}
		// fmt.Println("here")
	}()

	message, err := consumer.Receive(context.Background())

	if err == nil {
		fmt.Printf("consumed from topic %s at offset %v: "+
			string(message.Payload()), message.Topic())

		response := getModelResponse(message)
		writeRendezvousResponse(message, response)
	}
}

func getModelResponse(message pulsar.Message) []byte {
	var messageValue map[string]interface{}
	if err := json.Unmarshal(message.Payload(), &messageValue); err != nil {
		log.Fatal(err)
	}

	delete(messageValue, "id")

	requestBody, _ := json.Marshal(messageValue)

	modelEndpoint := os.Getenv("MODEL_ENDPOINT")

	response, err := http.Post(modelEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	return body
}

func writeRendezvousResponse(message pulsar.Message, body []byte) {
	var messageValue map[string]interface{}
	if err := json.Unmarshal(message.Payload(), &messageValue); err != nil {
		log.Fatal(err)
	}
	socketAddress := fmt.Sprintf("/sockets/%s.sock", messageValue["id"])

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

	response := fmt.Sprintf("{ \"results\": %s }", string(body))

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

	topic := config.TopicName("model-input")

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
