package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildUpdate handles discord.GatewayEventTypeGuildUpdate
type gatewayHandlerGuildUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildUpdate) New() interface{} {
	return &discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	guild := *v.(*discord.Guild)

	oldGuild := bot.Caches.Guilds().GetCopy(guild.ID)

	bot.EventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
			Guild:        bot.EntityBuilder.CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldGuild,
	})

}
