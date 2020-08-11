package main_test

import (
	"encoding/json"
	"github.com/concurrent-ai/rendezvous/src/shared/domain"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/concurrent-ai/rendezvous/src/model-executor"
	messagingMocks "github.com/concurrent-ai/rendezvous/src/shared/messaging/mocks"
)

func createFakeServer(response []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(response)
	}))
}

func getSentRendezvousMessage(mockProducer *messagingMocks.Producer) *domain.RendezvousMessage {
	var rendezvousMessage *domain.RendezvousMessage
	_ = json.Unmarshal(mockProducer.Calls[0].Arguments[0].([]byte), &rendezvousMessage)
	return rendezvousMessage
}

var _ = Describe("ModelExecutor", func() {
	Context("HandleNextMessage", func() {
		var (
			testMessage  *domain.RendezvousMessage
			mockConsumer *messagingMocks.Consumer
			mockProducer *messagingMocks.Producer
			fakeServer   *httptest.Server
			config       *Config
		)

		BeforeEach(func() {
			testMessage = &domain.RendezvousMessage{
				ID:          "TestID",
				RequestData: "test request",
			}

			testMessageBytes, _ := json.Marshal(testMessage)
			mockConsumer = &messagingMocks.Consumer{}
			mockConsumer.On("Receive").Return(testMessageBytes, nil)

			fakeModelResponse := []byte("test response")
			fakeServer = createFakeServer(fakeModelResponse)

			mockProducer = &messagingMocks.Producer{}
			mockProducer.On("Send", mock.Anything).Return(nil)

			config = &Config{
				OrganizationID: "TestOrg",
				ServiceID:      "TestService",
				ModelID:        "TestModel",
				ModelEndpoint:  fakeServer.URL,
				PulsarURL:      "pulsar://test:6650",
			}
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		It("should publish the received rendezvous message with the model id", func() {
			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.ResponseModelID).To(Equal("TestModel"))
		})

		It("should publish the received rendezvous message with the model response", func() {
			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.ResponseData).To(Equal("test response"))
		})

		It("should publish the received rendezvous message with a model request start event", func() {
			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.Events).To(Not(BeNil()))
			Expect(rendezvousMessage.Events.ModelRequestStart).To(Not(BeZero()))
		})

		It("should publish the received rendezvous message with a model request stop event", func() {
			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.Events).To(Not(BeNil()))
			Expect(rendezvousMessage.Events.ModelRequestStop).To(Not(BeZero()))
		})
	})
})
