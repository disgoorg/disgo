package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildUpdateHandler) New() interface{} {
	return &api.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	newGuild, ok := i.(*api.Guild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(newGuild.ID)
	if oldGuild != nil {
		oldGuild = &*oldGuild
	}
	newGuild = disgo.EntityBuilder().CreateGuild(newGuild, true)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		GuildID:      newGuild.ID,
	}

	eventManager.Dispatch(genericGuildEvent)
	eventManager.Dispatch(events.GuildUpdateEvent{
		GenericGuildEvent: genericGuildEvent,
		NewGuild:          newGuild,
		OldGuild:          oldGuild,
	})

}
