package domain

import (
	"time"
)

// RendezvousMessage : A struct that represents a single rendezvous request and response
type RendezvousMessage struct {
	ID              string            `json:"id"`
	RequestData     string            `json:"requestData"`
	ResponseModelID string            `json:"responseModelId,omitempty"`
	ResponseData    string            `json:"responseData,omitempty"`
	Events          *RendezvousEvents `json:"events"`
}

// RendezvousEvents : A struct that represents rendezvous message events
type RendezvousEvents struct {
	RequestTimestamp  time.Time `json:"requestReceived,omitempty"`
	ModelRequestStart time.Time `json:"modelRequestStart,omitempty"`
	ModelRequestStop  time.Time `json:"modelRequestStop,omitempty"`
}

// SetRequestTimestamp : Set the request timestamp event for a rendezvous message
func (message *RendezvousMessage) SetRequestTimestamp(timestamp time.Time) {
	if message.Events == nil {
		message.Events = &RendezvousEvents{}
	}
	message.Events.RequestTimestamp = timestamp
}

// SetModelRequestStart : Set the request timestamp event for a rendezvous message
func (message *RendezvousMessage) SetModelRequestStart(timestamp time.Time) {
	if message.Events == nil {
		message.Events = &RendezvousEvents{}
	}
	message.Events.ModelRequestStart = timestamp
}

// SetModelRequestStop : Set the request timestamp event for a rendezvous message
func (message *RendezvousMessage) SetModelRequestStop(timestamp time.Time) {
	if message.Events == nil {
		message.Events = &RendezvousEvents{}
	}
	message.Events.ModelRequestStop = timestamp
}
