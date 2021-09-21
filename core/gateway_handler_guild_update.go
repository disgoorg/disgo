package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildUpdate handles core.GuildUpdateGatewayEvent
type gatewayHandlerGuildUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildUpdate) New() interface{} {
	return &discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.Guild)

	oldCoreGuild := bot.Caches.GuildCache().GetCopy(guild.ID)

	bot.EventManager.Dispatch(&GuildUpdateEvent{
		GenericGuildEvent: &GenericGuildEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			Guild:        bot.EntityBuilder.CreateGuild(guild, CacheStrategyYes),
		},
		OldGuild: oldCoreGuild,
	})

}
