package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildRoleCreate struct{}

func (h *gatewayHandlerGuildRoleCreate) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleCreate
}

func (h *gatewayHandlerGuildRoleCreate) New() any {
	return &gateway.EventGuildRoleCreate{}
}

func (h *gatewayHandlerGuildRoleCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildRoleCreate)

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
