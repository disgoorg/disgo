package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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
	payload := *v.(*discord.UnavailableGuild)

	guild := bot.Caches.Guilds().Get(payload.ID)

	if payload.Unavailable {
		bot.Caches.Guilds().SetUnavailable(payload.ID)
	}

	genericGuildEvent := &events2.GenericGuildEvent{
		GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
		Guild:        guild,
	}

	if payload.Unavailable {
		bot.EventManager.Dispatch(&events2.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
