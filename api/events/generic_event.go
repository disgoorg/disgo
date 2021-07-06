package events

import "github.com/DisgoOrg/disgo/api"

// NewGenericEvent constructs a new GenericEvent with the provided Disgo instance
func NewGenericEvent(disgo api.Disgo, sequenceNumber int) *GenericEvent {
	return &GenericEvent{disgo: disgo, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	disgo          api.Disgo
	sequenceNumber int
}

// Disgo returns the Disgo instance for this event
func (e GenericEvent) Disgo() api.Disgo {
	return e.disgo
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}
