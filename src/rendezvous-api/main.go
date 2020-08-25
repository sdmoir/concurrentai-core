package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/concurrentai/concurrentai-core/src/shared/domain"
	"github.com/concurrentai/concurrentai-core/src/shared/messaging"
	"github.com/concurrentai/concurrentai-core/src/shared/sockets"
)

// main : Runs the rendezvous API server
func main() {
	log.Println("Starting server")

	config := LoadConfig()

	client, closeClient := createPulsarClient(config)
	defer closeClient()

	producer, closeProducer := createPulsarProducer(client, config)
	defer closeProducer()

	controller := NewAPIController(producer)
	http.HandleFunc("/", controller.HandleRequest)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

// createPulsarClient : Create a Pulsar client
func createPulsarClient(config *Config) (messaging.Client, func()) {
	client, err := messaging.NewPulsarClient(config.PulsarURL)
	if err != nil {
		log.Fatal(err)
	}
	return client, func() {
		client.Close()
	}
}

// createPulsarProducer : Create a Pulsar producer
func createPulsarProducer(client messaging.Client, config *Config) (messaging.Producer, func()) {
	producer, err := client.CreateProducer(config.TopicName("model-request"))
	if err != nil {
		log.Fatal(err)
	}
	return producer, func() {
		producer.Close()
	}
}

// APIController : Controller for handling rendezvous API requests
type APIController struct {
	Producer messaging.Producer
	Listener sockets.Listener
}

// NewAPIController : Creates a new APIController with the given message producer
func NewAPIController(producer messaging.Producer) *APIController {
	return &APIController{
		Producer: producer,
		Listener: sockets.NewUnixListener(),
	}
}

// HandleRequest : Handle a rendezvous API request
func (controller *APIController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.Wrap(err, "error reading request body"))
		http.Error(w, "can't ready body", http.StatusBadRequest)
		return
	}

	rendezvousMessage := createRendezvousMessage(body)
	producerPayload, err := json.Marshal(rendezvousMessage)
	if err != nil {
		log.Println(errors.Wrap(err, "error encoding rendezvous reqeuest payload"))
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}

	err = controller.Producer.Send(producerPayload)
	if err != nil {
		log.Println(errors.Wrap(err, "error sending rendezvous request payload"))
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}

	socketAddress := fmt.Sprintf("/sockets/%s.sock", rendezvousMessage.ID)
	rendezvousResponse, err := controller.Listener.Read(socketAddress)
	if err != nil {
		log.Println(errors.Wrap(err, "error reading rendezvous response from socket"))
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rendezvousResponse)
}

// createRendezvousMessage : Create a rendezvous message with the given body
func createRendezvousMessage(body []byte) *domain.RendezvousMessage {
	return &domain.RendezvousMessage{
		ID:          uuid.New().String(),
		RequestData: string(body),
		Events: &domain.RendezvousEvents{
			RequestTimestamp: time.Now(),
		},
	}
}
