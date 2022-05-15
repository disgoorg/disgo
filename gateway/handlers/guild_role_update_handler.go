package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildRoleUpdate handles discord.GatewayEventTypeGuildRoleUpdate
type gatewayHandlerGuildRoleUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) New() any {
	return &discord.GatewayEventGuildRoleUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildRoleUpdate)

	oldRole, _ := client.Caches().Roles().Get(payload.GuildID, payload.Role.ID)
	client.Caches().Roles().Put(payload.GuildID, payload.Role.ID, payload.Role)

	client.EventManager().DispatchEvent(&events.RoleUpdateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         payload.Role,
		},
		OldRole: oldRole,
	})
}
