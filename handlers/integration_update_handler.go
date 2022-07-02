package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerIntegrationUpdate struct {}

func (h *gatewayHandlerIntegrationUpdate) EventType() gateway.EventType {
	return gateway.EventTypeIntegrationUpdate
}

func (h *gatewayHandlerIntegrationUpdate) New() any {
	return &gateway.EventIntegrationUpdate{}
}

func (h *gatewayHandlerIntegrationUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventIntegrationUpdate)

	client.EventManager().DispatchEvent(&events.IntegrationUpdate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			Integration:  payload.Integration,
		},
	})
}
