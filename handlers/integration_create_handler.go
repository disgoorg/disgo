package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerIntegrationCreate struct{}

func (h *gatewayHandlerIntegrationCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationCreate
}

func (h *gatewayHandlerIntegrationCreate) New() any {
	return &discord.GatewayEventIntegrationCreate{}
}

func (h *gatewayHandlerIntegrationCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventIntegrationCreate)

	client.EventManager().DispatchEvent(&events.IntegrationCreate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			Integration:  payload.Integration,
		},
	})
}
