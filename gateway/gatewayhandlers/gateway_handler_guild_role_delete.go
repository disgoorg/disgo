package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildRoleDelete handles discord.GatewayEventTypeGuildRoleDelete
type gatewayHandlerGuildRoleDelete struct{}

// EventType returns the discord.GatewayEventType
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

	role := bot.Caches.RoleCache().GetCopy(payload.GuildID, payload.RoleID)

	bot.Caches.RoleCache().Remove(payload.GuildID, payload.RoleID)

	bot.EventManager.Dispatch(&events.RoleDeleteEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.RoleID,
			Role:   role,
		},
	})
}
