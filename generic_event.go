package disgo

// GenericEvent the basic interface each event implement
type GenericEvent interface {}

// Event the base event structure
type Event struct {
	Disgo          Disgo
	ResponseNumber int
}
