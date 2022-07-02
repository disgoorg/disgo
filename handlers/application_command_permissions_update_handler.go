package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerApplicationCommandPermissionsUpdate struct{}

func (h *gatewayHandlerApplicationCommandPermissionsUpdate) EventType() gateway.EventType {
	return gateway.EventTypeApplicationCommandPermissionsUpdate
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
