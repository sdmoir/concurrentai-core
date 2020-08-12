package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/concurrent-ai/rendezvous/src/model-enricher"
	messagingMocks "github.com/concurrent-ai/rendezvous/src/shared/messaging/mocks"
)

var _ = Describe("ModelEnricher", func() {
	Context("HandleNextMessage", func() {
		var (
			testMessage  []byte
			mockConsumer *messagingMocks.Consumer
			mockProducer *messagingMocks.Producer
		)

		BeforeEach(func() {
			testMessage = []byte("test")
			mockConsumer = &messagingMocks.Consumer{}
			mockProducer = &messagingMocks.Producer{}
		})

		It("should forward rendezvous messages to the model-input topic", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessage, nil)
			mockProducer.On("Send", testMessage).Return(nil)

			// act
			err := HandleNextMessage(mockConsumer, mockProducer)

			// assert
			Expect(err).To(Not(HaveOccurred()))
			Expect(mockProducer.AssertExpectations(GinkgoT())).To(BeTrue())
		})

		It("should return an error if it fails to read a rendezvous message from the consumer", func() {
			// arrange
			mockConsumer.On("Receive").Return(nil, errors.New("read error"))

			// act
			err := HandleNextMessage(mockConsumer, mockProducer)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read rendezvous message from consumer: read error"))
		})

		It("should return an error if it fails to publish the rendezvous message", func() {
			// arrange
			mockConsumer.On("Receive").Return(testMessage, nil)
			mockProducer.On("Send", testMessage).Return(errors.New("send error"))

			// act
			err := HandleNextMessage(mockConsumer, mockProducer)

			// assert
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send rendezvous message: send error"))
		})
	})
})
