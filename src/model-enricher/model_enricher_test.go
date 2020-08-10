package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/concurrent-ai/rendezvous/src/model-enricher"
	messagingMocks "github.com/concurrent-ai/rendezvous/src/shared/messaging/mocks"
)

var _ = Describe("ModelEnricher", func() {
	Context("HandleNextMessage", func() {
		It("should forward messages to the model-input topic", func() {
			// arrange
			testMessage := []byte("test")
			mockConsumer := &messagingMocks.Consumer{}
			mockConsumer.On("Receive").Return(testMessage, nil)
			mockProducer := &messagingMocks.Producer{}
			mockProducer.On("Send", testMessage).Return(nil)

			// act
			HandleNextMessage(mockConsumer, mockProducer)

			// assert
			Expect(mockProducer.AssertExpectations(GinkgoT())).To(BeTrue())
		})
	})
})
