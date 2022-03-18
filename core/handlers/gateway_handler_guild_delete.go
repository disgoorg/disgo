package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildDelete handles discord.GatewayEventTypeGuildDelete
type gatewayHandlerGuildDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildDelete) New() interface{} {
	return &discord.UnavailableGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildDelete) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.UnavailableGuild)

	guild := bot.Caches.Guilds().Get(payload.ID)

	if payload.Unavailable {
		bot.Caches.Guilds().SetUnavailable(payload.ID)
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		Guild:        guild,
	}

	if payload.Unavailable {
		bot.EventManager().Dispatch(&events.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager().Dispatch(&events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
