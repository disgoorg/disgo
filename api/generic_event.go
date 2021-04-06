package api

// Event the basic interface each event implement
type Event interface {
	Disgo() Disgo
}

// NewEvent constructs a new GenericEvent with the provided Disgo instance
func NewEvent(disgo Disgo) GenericEvent {
	return GenericEvent{disgo: disgo}
}

// GenericEvent the base event structure
type GenericEvent struct {
	disgo Disgo
}

// Disgo returns the Disgo instance for this event
func (d GenericEvent) Disgo() Disgo {
	return d.disgo
}
