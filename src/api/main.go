package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var config = LoadConfig()

type rendezvousRequest struct {
	id     string
	events map[string]time.Time
	data   map[string]interface{}
}

func readRequestBody(request *http.Request) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(bytes, &body); err != nil {
		return nil, errors.Wrap(err, "failed to parse request body")
	}

	return body, nil
}

func publishPulsarRequest(requestID uuid.UUID, request *http.Request) {
	client := CreatePulsarClient(config)
	defer client.Close()

	fmt.Println("Topic: " + config.TopicName("model-request"))

	producer := CreatePulsarProducer(client, config.TopicName("model-request"))
	defer producer.Close()

	body, err := readRequestBody(request)
	if err != nil {
		log.Fatal(err)
	}

	rendezvousRequest := &rendezvousRequest{
		id: requestID.String(),
		events: map[string]time.Time{
			"requestReceived": time.Now(),
		},
		data: body,
	}

	payloadBytes, err := json.Marshal(rendezvousRequest)
	if err != nil {
		log.Fatal("failed to marshal rendezvous request:", err)
	}

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payloadBytes,
	})

	if err != nil {
		log.Fatal("failed to publish message:", err)
	}

	log.Println("published rendezvous request: " + string(payloadBytes))
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
	log.Println(data)
}

func main() {
	log.Println("Starting server")
	http.HandleFunc("/", apiResponse)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
