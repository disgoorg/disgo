package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildRoleDelete struct{}

func (h *gatewayHandlerGuildRoleDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleDelete
}

func (h *gatewayHandlerGuildRoleDelete) New() any {
	return &discord.GatewayEventGuildRoleDelete{}
}

func (h *gatewayHandlerGuildRoleDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildRoleDelete)

	role, _ := client.Caches().Roles().Remove(payload.GuildID, payload.RoleID)

	client.EventManager().DispatchEvent(&events.RoleDelete{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			RoleID:       payload.RoleID,
			Role:         role,
		},
	})
}
