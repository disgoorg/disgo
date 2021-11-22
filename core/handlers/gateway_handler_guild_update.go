package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.Guild)

	oldGuild := bot.Caches.Guilds().GetCopy(guild.ID)

	bot.EventManager.Dispatch(&events2.GuildUpdateEvent{
		GenericGuildEvent: &events2.GenericGuildEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			Guild:        bot.EntityBuilder.CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldGuild,
	})

}
