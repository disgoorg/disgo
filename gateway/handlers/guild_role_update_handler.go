package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type roleUpdateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// GuildRoleUpdateHandler handles api.GuildRoleUpdateGatewayEvent
type GuildRoleUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildRoleUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleUpdateHandler) New() interface{} {
	return &roleUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload, ok := v.(roleUpdateData)
	if !ok {
		return
	}

	oldRole := bot.Caches.RoleCache().GetCopy(payload.GuildID, payload.Role.ID)

	bot.EventManager.Dispatch(&events.RoleUpdateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   bot.EntityBuilder.CreateRole(payload.GuildID, payload.Role, core.CacheStrategyYes),
		},
		OldRole: oldRole,
	})
}
