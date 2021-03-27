package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

// GuildDeleteHandler handles api.GuildDeleteGatewayEvent
type GuildDeleteHandler struct{}

// Name returns the raw gateway event name
func (h GuildDeleteHandler) Name() string {
	return api.GuildDeleteGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildDeleteHandler) New() interface{} {
	return &api.Guild{}
}

// Handle handles the specific raw gateway event
func (h GuildDeleteHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}

	if guild.Unavailable {
		disgo.Cache().Guild(guild.ID).Unavailable = true
		eventManager.Dispatch(events.GuildUnavailableEvent{
			GenericGuildEvent: events.GenericGuildEvent{
				Event: api.Event{
					Disgo: disgo,
				},
				GuildID: guild.ID,
			},
		})
	} else {
		cachedGuild := disgo.Cache().Guild(guild.ID)
		disgo.Cache().UncacheGuild(guild.ID)

		genericGuildEvent := events.GenericGuildEvent{
			Event: api.Event{
				Disgo: disgo,
			},
			GuildID: guild.ID,
		}

		eventManager.Dispatch(genericGuildEvent)

		eventManager.Dispatch(events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
			Guild:             cachedGuild,
		})
	}
}
