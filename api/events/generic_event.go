package events

import "github.com/DisgoOrg/disgo/api"

// NewEvent constructs a new GenericEvent with the provided Disgo instance
func NewEvent(disgo api.Disgo, sequenceNumber int) GenericEvent {
	return GenericEvent{disgo: disgo, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	disgo          api.Disgo
	sequenceNumber int
}

// Disgo returns the Disgo instance for this event
func (d GenericEvent) Disgo() api.Disgo {
	return d.disgo
}

// SequenceNumber returns the sequence number of the gateway event
func (d GenericEvent) SequenceNumber() int {
	return d.sequenceNumber
}
