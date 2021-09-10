package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type guildBanRemovePayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// gatewayHandlerGuildBanRemove handles core.GatewayEventGuildBanRemove
type gatewayHandlerGuildBanRemove struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildBanRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanRemove) New() interface{} {
	return &guildBanRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
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
