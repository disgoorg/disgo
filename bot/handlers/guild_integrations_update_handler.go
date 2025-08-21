package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildIntegrationsUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildIntegrationsUpdate) {
	client.EventManager.DispatchEvent(&events.GuildIntegrationsUpdate{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		GuildID:      event.GuildID,
	})
}
