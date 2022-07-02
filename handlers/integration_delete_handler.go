package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerIntegrationDelete struct {}

func (h *gatewayHandlerIntegrationDelete) EventType() gateway.EventType {
	return gateway.EventTypeIntegrationDelete
}

func (h *gatewayHandlerIntegrationDelete) New() any {
	return &gateway.EventIntegrationDelete{}
}

func (h *gatewayHandlerIntegrationDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventIntegrationDelete)

	client.EventManager().DispatchEvent(&events.IntegrationDelete{
		GenericEvent:  events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:       payload.GuildID,
		ID:            payload.ID,
		ApplicationID: payload.ApplicationID,
	})
}
