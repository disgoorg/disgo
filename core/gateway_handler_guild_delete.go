package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GuildDeleteHandler handles core.GuildDeleteGatewayEvent
type GuildDeleteHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *GuildDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildDeleteHandler) New() interface{} {
	return &discord.UnavailableGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildDeleteHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.UnavailableGuild)

	if guild.Unavailable {
		coreGuild := bot.Caches.GuildCache().Get(guild.ID)
		if coreGuild != nil {
			coreGuild.Unavailable = true
		}
	}

	genericGuildEvent := &GenericGuildEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		Guild:        bot.Caches.GuildCache().GetCopy(guild.ID),
	}

	if guild.Unavailable {
		bot.EventManager.Dispatch(&GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.Caches.GuildCache().Remove(guild.ID)

		bot.EventManager.Dispatch(&GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
