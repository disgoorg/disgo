package events

import "github.com/DiscoOrg/disgo"

// GenericEvent the basic interface each event implement
type GenericEvent interface {
	Disgo() disgo.Disgo
	ResponseNumber() int
}

// Event the base event structure
type Event struct {
	disgo          disgo.Disgo
	responseNumber int
}

// Disgo returns the disgo instance for this event
func (e Event) Disgo() disgo.Disgo {
	return e.disgo
}

// ResponseNumber returns the event
func (e Event) ResponseNumber() int {
	return e.responseNumber
}
