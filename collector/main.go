package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func handleNextMessage(consumer *kafka.Consumer) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic occurred: %s", err)
		}
		// fmt.Println("here")
	}()

	message, err := consumer.ReadMessage(100 * time.Millisecond)
	if err == nil {
		fmt.Printf("consumed from topic %s [%d] at offset %v: "+
			string(message.Value), *message.TopicPartition.Topic,
			message.TopicPartition.Partition, message.TopicPartition.Offset)

		response := getModelResponse(message)
		writeRendezvousResponse(message, response)
	}
}

func getModelResponse(message *kafka.Message) []byte {
	var messageValue map[string]interface{}
	if err := json.Unmarshal(message.Value, &messageValue); err != nil {
		log.Fatal(err)
	}

	delete(messageValue, "id")

	requestBody, _ := json.Marshal(messageValue)

	response, err := http.Post("http://model:8080/invocations", "application/json", bytes.NewBuffer(requestBody))
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

func writeRendezvousResponse(message *kafka.Message, body []byte) {
	var messageValue map[string]interface{}
	if err := json.Unmarshal(message.Value, &messageValue); err != nil {
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
	brokers := os.Getenv("KAFKA_BROKERS")
	apiKey := os.Getenv("KAFKA_API_KEY")
	apiSecret := os.Getenv("KAFKA_API_SECRET")
	topic := os.Getenv("KAFKA_TOPIC")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"sasl.mechanisms":    "PLAIN",
		"security.protocol":  "SASL_SSL",
		"sasl.username":      apiKey,
		"sasl.password":      apiSecret,
		"session.timeout.ms": 6000,
		"group.id":           "collector",
		"auto.offset.reset":  "latest"})

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create consumer: %s", err))
	}

	fmt.Println("Created consumer")

	topics := []string{topic}
	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		log.Fatal("error subscribing to topic:", err)
	}

	defer consumer.Close()

	for {
		handleNextMessage(consumer)
	}
}
