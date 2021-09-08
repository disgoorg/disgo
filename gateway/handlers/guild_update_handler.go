package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildUpdateHandler) New() interface{} {
	return discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild, ok := v.(discord.Guild)
	if !ok {
		return
	}

	oldCoreGuild := bot.Caches.GuildCache().GetCopy(guild.ID)

	bot.EventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			Guild:        bot.EntityBuilder.CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldCoreGuild,
	})

}
