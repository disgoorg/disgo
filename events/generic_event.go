package events

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
)

// NewGenericEvent constructs a new GenericEvent with the provided Client instance
func NewGenericEvent(client bot.Client, sequenceNumber discord.GatewaySequence) *GenericEvent {
	return &GenericEvent{client: client, sequenceNumber: sequenceNumber}
}

// GenericEvent the base event structure
type GenericEvent struct {
	client         bot.Client
	sequenceNumber discord.GatewaySequence
}

func (e GenericEvent) Client() bot.Client {
	return e.client
}

// SequenceNumber returns the sequence number of the gateway event
func (e GenericEvent) SequenceNumber() discord.GatewaySequence {
	return e.sequenceNumber
}
