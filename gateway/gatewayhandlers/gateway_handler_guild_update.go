package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.Guild)

	oldCoreGuild := bot.Caches.GuildCache().GetCopy(guild.ID)

	bot.EventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			Guild:        bot.EntityBuilder.CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldCoreGuild,
	})

}
