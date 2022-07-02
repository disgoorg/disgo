package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerIntegrationCreate struct{}

func (h *gatewayHandlerIntegrationCreate) EventType() gateway.EventType {
	return gateway.EventTypeIntegrationCreate
}

func (h *gatewayHandlerIntegrationCreate) New() any {
	return &gateway.EventIntegrationCreate{}
}

func (h *gatewayHandlerIntegrationCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventIntegrationCreate)

	client.EventManager().DispatchEvent(&events.IntegrationCreate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			Integration:  payload.Integration,
		},
	})
}
