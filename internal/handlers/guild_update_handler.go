package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildUpdateHandler) Event() api.GatewayEventName {
	return api.GatewayEventGuildUpdate
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
		GenericEvent: events.NewEvent(disgo),
		GuildID:      guild.ID,
	}

	eventManager.Dispatch(genericGuildEvent)
	eventManager.Dispatch(events.GuildUpdateEvent{
		GenericGuildEvent: genericGuildEvent,
		Guild:             guild,
		OldGuild:          &oldGuild,
	})

}
