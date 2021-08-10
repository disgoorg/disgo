package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// NewGenericEvent constructs a new GenericEvent with the provided Disgo instance
func NewGenericEvent(disgo core.Disgo, sequenceNumber int) *GenericEvent {
	return &GenericEvent{disgo: disgo, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	disgo          core.Disgo
	sequenceNumber int
}

// Disgo returns the Disgo instance for this event
func (e GenericEvent) Disgo() core.Disgo {
	return e.disgo
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}
