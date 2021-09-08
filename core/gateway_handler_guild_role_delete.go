package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type roleDeleteData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	RoleID  discord.Snowflake `json:"role_id"`
}

// GuildRoleDeleteHandler handles api.GuildRoleDeleteGatewayEvent
type GuildRoleDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildRoleDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleDeleteHandler) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleDeleteHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*roleDeleteData)

	role := bot.Caches.RoleCache().GetCopy(payload.GuildID, payload.RoleID)

	bot.Caches.RoleCache().Remove(payload.GuildID, payload.RoleID)

	bot.EventManager.Dispatch(&RoleDeleteEvent{
		GenericRoleEvent: &GenericRoleEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.RoleID,
			Role:   role,
		},
	})
}
