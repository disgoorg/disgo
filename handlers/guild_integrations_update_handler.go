package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildIntegrationsUpdate struct{}

func (h *gatewayHandlerGuildIntegrationsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildIntegrationsUpdate
}

func (h *gatewayHandlerGuildIntegrationsUpdate) New() any {
	return &discord.GatewayEventGuildIntegrationsUpdate{}
}

func (h *gatewayHandlerGuildIntegrationsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildIntegrationsUpdate)

	client.EventManager().DispatchEvent(&events.GuildIntegrationsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
	})
}
