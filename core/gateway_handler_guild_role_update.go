package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleUpdate handles discord.GatewayEventTypeGuildRoleUpdate
type gatewayHandlerGuildRoleUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) New() interface{} {
	return &discord.GuildRoleUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildRoleUpdateGatewayEvent)

	oldRole := bot.Caches.RoleCache().GetCopy(payload.GuildID, payload.Role.ID)

	bot.EventManager.Dispatch(&RoleUpdateEvent{
		GenericRoleEvent: &GenericRoleEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   bot.EntityBuilder.CreateRole(payload.GuildID, payload.Role, CacheStrategyYes),
		},
		OldRole: oldRole,
	})
}
