package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerApplicationCommandPermissionsUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventApplicationCommandPermissionsUpdate) {
	client.EventManager.DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		Permissions:  event.ApplicationCommandPermissions,
	})
}
