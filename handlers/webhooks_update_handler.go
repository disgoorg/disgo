package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerWebhooksUpdate struct{}

func (h *gatewayHandlerWebhooksUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeWebhooksUpdate
}

func (h *gatewayHandlerWebhooksUpdate) New() any {
	return &discord.GatewayEventWebhooksUpdate{}
}

func (h *gatewayHandlerWebhooksUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventWebhooksUpdate)

	client.EventManager().DispatchEvent(&events.WebhooksUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildId:      payload.GuildID,
		ChannelID:    payload.ChannelID,
	})
}
