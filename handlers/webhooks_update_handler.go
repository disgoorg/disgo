package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerWebhooksUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventWebhooksUpdate) {
	client.EventManager().DispatchEvent(&events.WebhooksUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildId:      event.GuildID,
		ChannelID:    event.ChannelID,
	})
}
