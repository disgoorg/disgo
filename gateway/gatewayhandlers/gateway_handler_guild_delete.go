package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildDelete handles core.GuildDeleteGatewayEvent
type gatewayHandlerGuildDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildDelete) New() interface{} {
	return &discord.UnavailableGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.UnavailableGuild)

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
