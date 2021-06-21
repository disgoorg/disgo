package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildDeleteHandler handles api.GuildDeleteGatewayEvent
type GuildDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h GuildDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildDeleteHandler) New() interface{} {
	return &api.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*api.FullGuild)
	if !ok {
		return
	}

	guild := disgo.EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyNo)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Guild:        guild,
	}
	eventManager.Dispatch(genericGuildEvent)

	if guild.Unavailable {
		// set guild to unavailable for now
		disgo.Cache().Guild(guild.ID).Unavailable = true

		eventManager.Dispatch(events.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		disgo.Cache().UncacheGuild(guild.ID)

		eventManager.Dispatch(events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
