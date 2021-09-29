package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildBanRemove handles discord.GatewayEventTypeGuildBanRemove
type gatewayHandlerGuildBanRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildBanRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanRemove) New() interface{} {
	return &discord.GuildBanRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildBanRemoveGatewayEvent)

	bot.EventManager.Dispatch(&GuildUnbanEvent{
		GenericGuildEvent: &GenericGuildEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
		},
		User: bot.EntityBuilder.CreateUser(payload.User, CacheStrategyNo),
	})
}
