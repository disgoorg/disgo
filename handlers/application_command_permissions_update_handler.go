package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerApplicationCommandPermissionsUpdate struct{}

func (h *gatewayHandlerApplicationCommandPermissionsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeApplicationCommandPermissionsUpdate
}

func (h *gatewayHandlerApplicationCommandPermissionsUpdate) New() any {
	return &discord.ApplicationCommandPermissions{}
}

func (h *gatewayHandlerApplicationCommandPermissionsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	permissions := *v.(*discord.ApplicationCommandPermissions)

	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		Permissions:  permissions,
	})
}
