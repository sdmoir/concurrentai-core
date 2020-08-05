package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"io/ioutil"
	"log"
	"net/http"
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

func handleNextMessage(consumer pulsar.Consumer, producer pulsar.Producer, modelEndpoint string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
	}()

	message, err := consumer.Receive(context.Background())

	if err == nil {
		log.Printf("consumed from topic " + message.Topic() + ": " + string(message.Payload()))

		response := getModelResponse(message, modelEndpoint)
		payload, _ := json.Marshal(response)

		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: payload,
		})

		log.Print("published to topic " + producer.Topic())
	}
}

func getModelResponse(message pulsar.Message, modelEndpoint string) map[string]interface{} {
	var messageValue map[string]interface{}
	if err := json.Unmarshal(message.Payload(), &messageValue); err != nil {
		log.Fatal(err)
	}

	id := messageValue["id"]
	delete(messageValue, "id")

	requestBody, _ := json.Marshal(messageValue)

	response, err := http.Post(modelEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"id":       id,
		"response": string(body),
	}
}

func main() {
	config := LoadConfig()

	client := createPulsarClient(config)
	defer client.Close()

	consumer := createPulsarConsumer(client, config.TopicName("model-input"))
	defer consumer.Close()

	producer := createPulsarProducer(client, config.TopicName("model-response"))
	defer producer.Close()

	for {
		handleNextMessage(consumer, producer, config.modelEndpoint)
	}
}
