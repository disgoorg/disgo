package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerResumed struct{}

func (h *gatewayHandlerResumed) EventType() gateway.EventType {
	return gateway.EventTypeResumed
}

func (h *gatewayHandlerResumed) New() any {
	return nil
}

func (h *gatewayHandlerResumed) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, _ any) {
	client.EventManager().DispatchEvent(&events.Resumed{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
	})
}
