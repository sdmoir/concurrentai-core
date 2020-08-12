package main_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	. "github.com/concurrent-ai/rendezvous/src/model-executor"
	"github.com/concurrent-ai/rendezvous/src/shared/domain"
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
			testMessage      *domain.RendezvousMessage
			testMessageBytes []byte
			mockConsumer     *messagingMocks.Consumer
			mockProducer     *messagingMocks.Producer
			fakeServer       *httptest.Server
			config           *Config
		)

		BeforeEach(func() {
			testMessage = &domain.RendezvousMessage{
				ID:          "test message",
				RequestData: "test request",
			}

			testMessageBytes, _ = json.Marshal(testMessage)
			mockConsumer = &messagingMocks.Consumer{}
			mockProducer = &messagingMocks.Producer{}

			fakeModelResponse := []byte(`["test response"]`)
			fakeServer = createFakeServer(fakeModelResponse)

			config = &Config{
				OrganizationID: "test org",
				ServiceID:      "test service",
				ModelID:        "test model",
				ModelEndpoint:  fakeServer.URL,
				PulsarURL:      "pulsar://test:6650",
			}
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		It("should publish the received rendezvous message with the model id", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockProducer.On("Send", mock.Anything).Return(nil)

			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.ResponseModelID).To(Equal("test model"))
		})

		It("should publish the received rendezvous message with the model response", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockProducer.On("Send", mock.Anything).Return(nil)

			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.ResponseData).To(Equal(`{ "results": ["test response"] }`))
		})

		It("should publish the received rendezvous message with a model request start event", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockProducer.On("Send", mock.Anything).Return(nil)

			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.Events).To(Not(BeNil()))
			Expect(rendezvousMessage.Events.ModelRequestStart).To(Not(BeZero()))
		})

		It("should publish the received rendezvous message with a model request stop event", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockProducer.On("Send", mock.Anything).Return(nil)

			// act
			HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.Events).To(Not(BeNil()))
			Expect(rendezvousMessage.Events.ModelRequestStop).To(Not(BeZero()))
		})

		It("should return an error if it fails to read a rendezvous message from the consumer", func() {
			// arrange
			mockConsumer.On("Receive").Return(nil, errors.New("read error"))

			// act
			err := HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read rendezvous message from consumer: read error"))
		})

		It("should return an error if it fails to parse a rendezvous message", func() {
			// arrange
			mockConsumer.On("Receive").Return([]byte{}, nil)

			// act
			err := HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to parse rendezvous message: unexpected end of JSON input"))
		})

		It("should return an error if it fails to get the model response", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			config.ModelEndpoint = ""

			// act
			err := HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`failed to get model response: error calling model endpoint: Post "": unsupported protocol scheme ""`))
		})

		It("should return an error if it fails to send the rendezvous message with the model response", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockProducer.On("Send", mock.Anything).Return(errors.New("send error"))

			// act
			err := HandleNextMessage(mockConsumer, mockProducer, config)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send rendezvous message with model response: send error"))
		})
	})
})
