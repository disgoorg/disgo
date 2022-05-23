package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildRoleUpdate struct{}

func (h *gatewayHandlerGuildRoleUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

func (h *gatewayHandlerGuildRoleUpdate) New() any {
	return &discord.GatewayEventGuildRoleUpdate{}
}

func (h *gatewayHandlerGuildRoleUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildRoleUpdate)

	oldRole, _ := client.Caches().Roles().Get(payload.GuildID, payload.Role.ID)
	client.Caches().Roles().Put(payload.GuildID, payload.Role.ID, payload.Role)

	client.EventManager().DispatchEvent(&events.RoleUpdate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         payload.Role,
		},
		OldRole: oldRole,
	})
}
