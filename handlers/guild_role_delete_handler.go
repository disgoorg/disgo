package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildRoleDelete struct {}

func (h *gatewayHandlerGuildRoleDelete) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleDelete
}

func (h *gatewayHandlerGuildRoleDelete) New() any {
	return &gateway.EventGuildRoleDelete{}
}

func (h *gatewayHandlerGuildRoleDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildRoleDelete)

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
