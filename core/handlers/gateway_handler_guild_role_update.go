package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleUpdate handles discord.GatewayEventTypeGuildRoleUpdate
type gatewayHandlerGuildRoleUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) New() any {
	return &discord.GuildRoleUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildRoleUpdateGatewayEvent)

	oldRole, _ := bot.Caches().Roles().Get(payload.GuildID, payload.Role.ID)
	bot.Caches().Roles().Put(payload.GuildID, payload.Role.ID, payload.Role)

	bot.EventManager().Dispatch(&events.RoleUpdateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         payload.Role,
		},
		OldRole: oldRole,
	})
}
