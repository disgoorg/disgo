package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildDeleteHandler handles api.GuildDeleteGatewayEvent
type GuildDeleteHandler struct{}

// EventType returns the api.GatewayEventType
func (h *GuildDeleteHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildDeleteHandler) New() interface{} {
	return &discord.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*discord.FullGuild)
	if !ok {
		return
	}

	guild := disgo.EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyNo)

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		Guild:        guild,
	}

	if guild.Unavailable {
		// set guild to unavailable for now
		g := disgo.Cache().Guild(guild.ID)
		if g != nil {
			g.Unavailable = true
		}

		eventManager.Dispatch(&events.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		disgo.Cache().UncacheGuild(guild.ID)

		eventManager.Dispatch(&events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
