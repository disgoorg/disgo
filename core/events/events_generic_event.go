package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// NewGenericEvent constructs a new GenericEvent with the provided Bot instance
func NewGenericEvent(bot core.Bot, sequenceNumber int) *GenericEvent {
	return &GenericEvent{bot: bot, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	bot            *core.Bot
	sequenceNumber int
}

func (e GenericEvent) Bot() *core.Bot {
	return e.bot
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}
