package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
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

func handleNextMessage(consumer pulsar.Consumer, producer pulsar.Producer, config Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s", err)
		}
	}()

	message, err := consumer.Receive(context.Background())

	if err == nil {
		log.Printf("consumed from topic " + message.Topic() + ": " + string(message.Payload()))

		request, err := parseModelRequest(message)
		if err != nil {
			log.Println(err)
			return
		}

		response, err := getModelResponse(request, config)
		if err != nil {
			log.Println(err)
			return
		}

		payload, _ := json.Marshal(response)
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{Payload: payload})
		if err != nil {
			log.Println(errors.Wrap(err, "failed to publish message for model response"))
			return
		}

		log.Print("published to topic " + producer.Topic())
	}
}

type rendezvousModelRequest struct {
	id     string
	events map[string]time.Time
	body   map[string]interface{}
}

type rendezvousModelResponse struct {
	id            string
	events        map[string]interface{}
	body          map[string]interface{}
	modelID       string
	modelResponse string
}

func parseModelRequest(message pulsar.Message) (*rendezvousModelRequest, error) {
	var request *rendezvousModelRequest
	if err := json.Unmarshal(message.Payload(), &request); err != nil {
		return nil, errors.Wrap(err, "error parsing rendezvous request from message")
	}
	return request, nil
}

func getModelResponse(request *rendezvousModelRequest, config Config) (*rendezvousModelResponse, error) {
	requestBody, _ := json.Marshal(request.body)

	startTime := time.Now()

	response, err := http.Post(config.modelEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.Wrap(err, "error calling model endpoint")
	}

	defer response.Body.Close()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading model response")
	}

	rendezvousResponseEvents := map[string]interface{}{
		"modelRequestStart":                startTime,
		"modelRequestStop":                 endTime,
		"modelRequestDurationMilliseconds": duration.Milliseconds(),
	}

	for key, value := range request.events {
		rendezvousResponseEvents[key] = value
	}

	rendezvousResponse := &rendezvousModelResponse{
		id:            request.id,
		events:        rendezvousResponseEvents,
		body:          request.body,
		modelID:       config.modelID,
		modelResponse: string(body),
	}

	return rendezvousResponse, nil
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
		handleNextMessage(consumer, producer, config)
	}
}
