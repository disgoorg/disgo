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
	return &api.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}

	guild = disgo.EntityBuilder().CreateGuild(guild, api.CacheStrategyNo)

	if guild.Unavailable {
		// set guild to unavail for now
		disgo.Cache().Guild(guild.ID).Unavailable = true
		eventManager.Dispatch(events.GuildUnavailableEvent{
			GenericGuildEvent: events.GenericGuildEvent{
				GenericEvent: events.NewEvent(disgo, sequenceNumber),
				GuildID:      guild.ID,
			},
		})
	} else {
		guild = disgo.Cache().Guild(guild.ID)
		disgo.Cache().UncacheGuild(guild.ID)

		genericGuildEvent := events.GenericGuildEvent{
			GenericEvent: events.NewEvent(disgo, sequenceNumber),
			GuildID:      guild.ID,
		}
		eventManager.Dispatch(genericGuildEvent)

		eventManager.Dispatch(events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
			Guild:             guild,
		})
	}
}
