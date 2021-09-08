package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildDeleteHandler handles api.GuildDeleteGatewayEvent
type GuildDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildDeleteHandler) New() interface{} {
	return discord.UnavailableGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildDeleteHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild, ok := v.(discord.UnavailableGuild)
	if !ok {
		return
	}

	if guild.Unavailable {
		coreGuild := bot.Caches.GuildCache().Get(guild.ID)
		if coreGuild != nil {
			coreGuild.Unavailable = true
		}
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		Guild:        bot.Caches.GuildCache().GetCopy(guild.ID),
	}

	if guild.Unavailable {
		bot.EventManager.Dispatch(&events.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.Caches.GuildCache().Remove(guild.ID)

		bot.EventManager.Dispatch(&events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
