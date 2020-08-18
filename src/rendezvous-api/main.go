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

// APIController : Controller for handling rendezvous API requests
type APIController struct {
	Producer messaging.Producer
	Listener sockets.Listener
}

// HandleRequest : Handle a rendezvous API request
func (controller *APIController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.Wrap(err, "error reading request body"))
		http.Error(w, "can't ready body", http.StatusBadRequest)
		return
	}

	rendezvousMessage := &domain.RendezvousMessage{
		ID:          uuid.New().String(),
		RequestData: string(body),
		Events: &domain.RendezvousEvents{
			RequestTimestamp: time.Now(),
		},
	}

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
	log.Println("Waiting for response at " + socketAddress)
	rendezvousResponse, err := controller.Listener.Read(socketAddress)
	if err != nil {
		log.Println(errors.Wrap(err, "error reading rendezvous response from socket"))
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rendezvousResponse)
}

func main() {
	log.Println("Starting server")

	config := LoadConfig()

	pulsarClient, err := messaging.NewPulsarClient(config.PulsarURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pulsarClient.Close()

	pulsarProducer, err := pulsarClient.CreateProducer(config.TopicName("model-request"))
	if err != nil {
		log.Fatal(err)
	}
	defer pulsarProducer.Close()

	controller := &APIController{
		Producer: pulsarProducer,
		Listener: sockets.NewUnixListener(),
	}

	http.HandleFunc("/", controller.HandleRequest)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
