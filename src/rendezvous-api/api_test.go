package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	. "github.com/concurrentai/concurrentai-core/src/rendezvous-api"
	"github.com/concurrentai/concurrentai-core/src/shared/domain"
	messagingMocks "github.com/concurrentai/concurrentai-core/src/shared/messaging/mocks"
	socketMocks "github.com/concurrentai/concurrentai-core/src/shared/sockets/mocks"
)

type badReader int

func (br badReader) Read(payload []byte) (int, error) {
	return 0, errors.New("some error while reading body")
}

func getSentRendezvousMessage(mockProducer *messagingMocks.Producer) *domain.RendezvousMessage {
	var rendezvousMessage *domain.RendezvousMessage
	_ = json.Unmarshal(mockProducer.Calls[0].Arguments[0].([]byte), &rendezvousMessage)
	return rendezvousMessage
}

var _ = Describe("APIController", func() {
	Context("HandleRequest", func() {
		const (
			method = "GET"
			url    = "http://concurrent.ai/test"
		)

		var (
			mockProducer *messagingMocks.Producer
			mockListener *socketMocks.Listener
			controller   *APIController
			body         string
			req          *http.Request
			w            *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			mockProducer = &messagingMocks.Producer{}
			mockListener = &socketMocks.Listener{}
			controller = &APIController{
				Producer: mockProducer,
				Listener: mockListener,
			}

			body = `{
				"columns": ["test"],
				"data": [["1"]]
			}`
			req = httptest.NewRequest(method, url, strings.NewReader(body))
			w = httptest.NewRecorder()
		})

		It("should return a 400 status when the request body can't be read", func() {
			// arrange
			req = httptest.NewRequest(method, url, badReader(0))

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should publish a rendezvous message with id to pulsar", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.ID).To(Not(BeEmpty()))
		})

		It("should publish a rendezvous message with request data to pulsar", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.RequestData).To(Equal(body))
		})

		It("should publish a rendezvous message with request timestamp event to pulsar", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			Expect(rendezvousMessage.Events.RequestTimestamp).To(Not(BeZero()))
		})

		It("should return a 500 status if publishing the rendezvous message fails", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(errors.New("error publishing rendezvous requeset"))

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should listen to the correct socket for the rendezvous response", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			rendezvousMessage := getSentRendezvousMessage(mockProducer)
			expectedSocketAddress := fmt.Sprint("/sockets/" + rendezvousMessage.ID + ".sock")
			mockListener.AssertCalled(GinkgoT(), "Read", expectedSocketAddress)
		})

		It("should write response data received from the socket listener", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)
			Expect(responseBody).To(Equal([]byte(`{ "data": "test" }`)))
		})

		It("should return a 500 status when an error occurs reading from the socket listener", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return(nil, errors.New("error reading from socket"))

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should return a 200 status when successful", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("should set the content type to application/json in the response", func() {
			// arrange
			mockProducer.On("Send", mock.Anything).Return(nil).Once()
			mockListener.On("Read", mock.Anything).Return([]byte(`{ "data": "test" }`), nil)

			// act
			controller.HandleRequest(w, req)

			// assert
			response := w.Result()
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})
})
