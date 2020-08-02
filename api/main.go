package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var config = LoadConfig()

func publishPulsarRequest(requestID uuid.UUID, request *http.Request) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: config.pulsarURL,
	})
	if err != nil {
		log.Fatal("failed to create pulsar client", err)
	}
	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: TopicName(config, "model-request"),
	})
	if err != nil {
		log.Fatal("failed to create pulsar producer", err)
	}
	defer producer.Close()

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

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payloadBytes,
	})

	if err != nil {
		log.Fatal("failed to publish message:", err)
	}
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
	publishPulsarRequest(requestID, r)
	data := waitForRendezvousResponse(requestID)
	writeResponseData(w, data)
	fmt.Println(data)
}

func main() {
	fmt.Println("Starting server")
	http.HandleFunc("/", apiResponse)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
