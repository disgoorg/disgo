package events

import "github.com/DisgoOrg/disgo/api"

// NewEvent constructs a new GenericEvent with the provided Disgo instance
func NewEvent(disgo api.Disgo) GenericEvent {
	return GenericEvent{disgo: disgo}
}

// GenericEvent the base event structure
type GenericEvent struct {
	disgo api.Disgo
}

// Disgo returns the Disgo instance for this event
func (d GenericEvent) Disgo() api.Disgo {
	return d.disgo
}
