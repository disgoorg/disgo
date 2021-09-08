package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type roleUpdateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// GuildRoleUpdateHandler handles core.GuildRoleUpdateGatewayEvent
type GuildRoleUpdateHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *GuildRoleUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleUpdateHandler) New() interface{} {
	return &roleUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleUpdateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*roleUpdateData)

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
