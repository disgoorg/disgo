package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildRoleDelete handles discord.GatewayEventTypeGuildRoleDelete
type gatewayHandlerGuildRoleDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleDelete) New() any {
	return &discord.GatewayEventGuildRoleDelete{}
}

// HandleGatewayEvent handles the specific raw gateway event
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
