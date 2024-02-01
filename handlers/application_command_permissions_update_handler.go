package handlers

import (
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
)

func gatewayHandlerApplicationCommandPermissionsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventApplicationCommandPermissionsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		Permissions:  event.ApplicationCommandPermissions,
	})
}
