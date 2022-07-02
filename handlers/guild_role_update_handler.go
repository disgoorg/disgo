package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildRoleUpdate struct{}

func (h *gatewayHandlerGuildRoleUpdate) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleUpdate
}

func (h *gatewayHandlerGuildRoleUpdate) New() any {
	return &gateway.EventGuildRoleUpdate{}
}

func (h *gatewayHandlerGuildRoleUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildRoleUpdate)

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
