package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type GuildUpdateHandler struct{}

// Name returns the raw gateway event name
func (h GuildUpdateHandler) Name() string {
	return api.GuildUpdateGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildUpdateHandler) New() interface{} {
	return &api.Guild{}
}

// Handle handles the specific raw gateway event
func (h GuildUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}

	oldGuild := *disgo.Cache().Guild(guild.ID)
	disgo.Cache().CacheGuild(guild)

	genericGuildEvent := events.GenericGuildEvent{
		Event:   api.Event{
			Disgo: disgo,
		},
		GuildID: guild.ID,
	}

	eventManager.Dispatch(genericGuildEvent)
	eventManager.Dispatch(events.GuildUpdateEvent{
		GenericGuildEvent: genericGuildEvent,
		Guild:             guild,
		OldGuild:          &oldGuild,
	})

}
