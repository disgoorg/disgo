package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerResumed(client bot.Client, sequenceNumber int, shardID int, _ gateway.EventData) {
	client.EventManager().DispatchEvent(&events.Resumed{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
	})
}
