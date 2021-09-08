package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type roleCreateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// GuildRoleCreateHandler handles api.GuildRoleCreateGatewayEvent
type GuildRoleCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildRoleCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleCreateHandler) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleCreateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload, ok := v.(roleCreateData)
	if !ok {
		return
	}

	bot.EventManager.Dispatch(&events.RoleCreateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   bot.EntityBuilder.CreateRole(payload.GuildID, payload.Role, core.CacheStrategyYes),
		},
	})
}
