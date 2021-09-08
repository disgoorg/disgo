package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type guildBanRemovePayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// GuildBanRemoveHandler handles core.GatewayEventGuildBanRemove
type GuildBanRemoveHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *GuildBanRemoveHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanRemoveHandler) New() interface{} {
	return &guildBanRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanRemoveHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*guildBanRemovePayload)

	bot.EventManager.Dispatch(&GuildUnbanEvent{
		GenericGuildEvent: &GenericGuildEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
		},
		User: bot.EntityBuilder.CreateUser(payload.User, CacheStrategyNo),
	})
}
