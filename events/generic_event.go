package events

import (
	"github.com/disgoorg/disgo/bot"
)

// NewGenericEvent constructs a new GenericEvent with the provided Client instance
func NewGenericEvent(client bot.Client, sequenceNumber int) *GenericEvent {
	return &GenericEvent{client: client, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	client         bot.Client
	sequenceNumber int
}

func (e GenericEvent) Client() bot.Client {
	return e.client
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}
