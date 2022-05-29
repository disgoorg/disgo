package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerIntegrationDelete struct{}

func (h *gatewayHandlerIntegrationDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationDelete
}

func (h *gatewayHandlerIntegrationDelete) New() any {
	return &discord.GatewayEventIntegrationDelete{}
}

func (h *gatewayHandlerIntegrationDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventIntegrationDelete)

	client.EventManager().DispatchEvent(&events.IntegrationDelete{
		GenericEvent:  events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:       payload.GuildID,
		ID:            payload.ID,
		ApplicationID: payload.ApplicationID,
	})
}
