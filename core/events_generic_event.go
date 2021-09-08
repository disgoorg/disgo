package core

// NewGenericEvent constructs a new GenericEvent with the provided Bot instance
func NewGenericEvent(bot *Bot, sequenceNumber int) *GenericEvent {
	return &GenericEvent{bot: bot, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	bot            *Bot
	sequenceNumber int
}

func (e GenericEvent) Bot() *Bot {
	return e.bot
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}
