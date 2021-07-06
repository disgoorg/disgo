package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h *GuildUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildUpdateHandler) New() interface{} {
	return &api.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*api.FullGuild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(fullGuild.ID)
	if oldGuild != nil {
		oldGuild = &*oldGuild
	}
	guild := disgo.EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyYes)

	eventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Guild:        guild,
		},
		OldGuild: oldGuild,
	})

}
