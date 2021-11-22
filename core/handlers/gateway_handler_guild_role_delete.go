package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleDelete handles core.GuildRoleDeleteGatewayEvent
type gatewayHandlerGuildRoleDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildRoleDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleDelete) New() interface{} {
	return &discord.GuildRoleDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildRoleDeleteGatewayEvent)

	role := bot.Caches.Roles().GetCopy(payload.GuildID, payload.RoleID)

	bot.Caches.Roles().Remove(payload.GuildID, payload.RoleID)

	bot.EventManager.Dispatch(&events2.RoleDeleteEvent{
		GenericRoleEvent: &events2.GenericRoleEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			RoleID:       payload.RoleID,
			Role:         role,
		},
	})
}
