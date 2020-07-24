package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func publishRequest(requestID uuid.UUID, request *http.Request) {
	brokers := os.Getenv("KAFKA_BROKERS")
	apiKey := os.Getenv("KAFKA_API_KEY")
	apiSecret := os.Getenv("KAFKA_API_SECRET")
	topic := os.Getenv("KAFKA_TOPIC")

	fmt.Printf("%s", brokers)
	fmt.Printf("%s", apiKey)
	fmt.Printf("%s", apiSecret)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     apiKey,
		"sasl.password":     apiSecret})

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create producer: %s", err))
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal("failed to read request body:", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Fatal("failed to parse request body:", err)
	}

	payload["id"] = requestID.String()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal request body:", err)
	}

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
		Value: payloadBytes}, nil)

	// Wait for delivery report
	e := <-producer.Events()

	switch e.(type) {
	default:
		fmt.Printf("produced message: %s", payload["id"])
	case kafka.Error:
		log.Fatal("kafka error:", e)
	}

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		fmt.Printf("failed to deliver message: %v\n",
			message.TopicPartition)
	} else {
		fmt.Printf("delivered to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic,
			message.TopicPartition.Partition,
			message.TopicPartition.Offset)
	}

	producer.Close()
}

func waitForRendezvousResponse(requestID uuid.UUID) []byte {
	socketAddress := fmt.Sprintf("/sockets/%s.sock", requestID)
	data := waitForSocketData(socketAddress)
	return data
}

func waitForSocketData(socketAddress string) []byte {
	listener, err := net.Listen("unix", socketAddress)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	defer conn.Close()

	data, error := ioutil.ReadAll(conn)
	if error != nil {
		log.Fatal("read error:", error)
	}

	return data
}

func writeResponseData(response http.ResponseWriter, data []byte) {
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(data))
}

func apiResponse(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()
	publishRequest(requestID, r)
	data := waitForRendezvousResponse(requestID)
	writeResponseData(w, data)
	fmt.Println(data)
}

func main() {
	fmt.Println("Starting server")
	http.HandleFunc("/", apiResponse)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
