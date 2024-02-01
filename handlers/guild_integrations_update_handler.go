package handlers

import (
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
)

func gatewayHandlerGuildIntegrationsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildIntegrationsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildIntegrationsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	})
}
