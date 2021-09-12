package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type roleCreateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// gatewayHandlerGuildRoleCreate handles core.GuildRoleCreateGatewayEvent
type gatewayHandlerGuildRoleCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildRoleCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleCreate) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleCreate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*roleCreateData)

	bot.EventManager.Dispatch(&RoleCreateEvent{
		GenericRoleEvent: &GenericRoleEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   bot.EntityBuilder.CreateRole(payload.GuildID, payload.Role, CacheStrategyYes),
		},
	})
}
