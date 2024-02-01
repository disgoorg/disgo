package handlers

import (
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
)

func gatewayHandlerWebhooksUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventWebhooksUpdate) {
	client.EventManager().DispatchEvent(&events.WebhooksUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildId:      event.GuildID,
		ChannelID:    event.ChannelID,
	})
}
