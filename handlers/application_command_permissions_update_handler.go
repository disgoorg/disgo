package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerApplicationCommandPermissionsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventApplicationCommandPermissionsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		Permissions:  event.ApplicationCommandPermissions,
	})
}
