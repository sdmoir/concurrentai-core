package main_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/concurrent-ai/rendezvous/src/collector"
	"github.com/concurrent-ai/rendezvous/src/shared/domain"
	messagingMocks "github.com/concurrent-ai/rendezvous/src/shared/messaging/mocks"
	socketMocks "github.com/concurrent-ai/rendezvous/src/shared/sockets/mocks"
)

var _ = Describe("Collector", func() {
	Context("HandleNextMessage", func() {
		var (
			testMessage       *domain.RendezvousMessage
			testMessageBytes  []byte
			fakeSocketAddress string
			mockConsumer      *messagingMocks.Consumer
			mockSocketWriter  *socketMocks.Writer
		)

		BeforeEach(func() {
			testMessage = &domain.RendezvousMessage{
				ID:              "test message",
				RequestData:     "test request",
				ResponseModelID: "test model",
				ResponseData:    "test response",
			}

			testMessageBytes, _ = json.Marshal(testMessage)
			mockConsumer = &messagingMocks.Consumer{}

			fakeSocketAddress = "/sockets/" + testMessage.ID + ".sock"
			mockSocketWriter = &socketMocks.Writer{}
		})

		It("should read a rendezvous message from the consumer", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockSocketWriter.On("Write", fakeSocketAddress, []byte(testMessage.ResponseData)).Return(nil)

			// act
			err := HandleNextMessage(mockConsumer, mockSocketWriter)

			// assert
			Expect(err).To(Not(HaveOccurred()))
			Expect(mockConsumer.AssertExpectations(GinkgoT())).To(BeTrue())
		})

		It("should return an error if it fails to read a rendezvous message from the consumer", func() {
			// arrange
			mockConsumer.On("Receive").Return(nil, errors.New("read error"))

			// act
			err := HandleNextMessage(mockConsumer, mockSocketWriter)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read rendezvous message from consumer: read error"))
		})

		It("should return an error if it fails to parse a rendezvous message", func() {
			// arrange
			mockConsumer.On("Receive").Return([]byte{}, nil)

			// act
			err := HandleNextMessage(mockConsumer, mockSocketWriter)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to parse rendezvous message: unexpected end of JSON input"))
		})

		It("should write the rendezvous message response data to the expected socket", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockSocketWriter.On("Write", fakeSocketAddress, []byte(testMessage.ResponseData)).Return(nil)

			// act
			err := HandleNextMessage(mockConsumer, mockSocketWriter)

			// assert
			Expect(err).To(Not(HaveOccurred()))
			Expect(mockSocketWriter.AssertExpectations(GinkgoT())).To(BeTrue())
		})

		It("should return an error if it fails to write the model response to the expected socket", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessageBytes, nil)
			mockSocketWriter.On("Write", fakeSocketAddress, []byte(testMessage.ResponseData)).Return(errors.New("write error"))

			// act
			err := HandleNextMessage(mockConsumer, mockSocketWriter)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to write rendezvous message response data to socket: write error"))
		})
	})
})
