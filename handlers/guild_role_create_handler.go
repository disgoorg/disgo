package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildRoleCreate struct{}

func (h *gatewayHandlerGuildRoleCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleCreate
}

func (h *gatewayHandlerGuildRoleCreate) New() any {
	return &discord.GatewayEventGuildRoleCreate{}
}

func (h *gatewayHandlerGuildRoleCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildRoleCreate)

	client.Caches().Roles().Put(payload.GuildID, payload.Role.ID, payload.Role)

	client.EventManager().DispatchEvent(&events.RoleCreate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         payload.Role,
		},
	})
}
