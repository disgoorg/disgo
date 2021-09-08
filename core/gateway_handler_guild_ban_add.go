package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type guildBanAddPayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// GuildBanAddHandler handles api.GatewayEventGuildBanAdd
type GuildBanAddHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildBanAddHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanAddHandler) New() interface{} {
	return &guildBanAddPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanAddHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*guildBanAddPayload)

	bot.EventManager.Dispatch(&GuildBanEvent{
		GenericGuildEvent: &GenericGuildEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
		},
		User: bot.EntityBuilder.CreateUser(payload.User, CacheStrategyNo),
	})
}
