package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildUpdateHandler) New() interface{} {
	return &discord.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*discord.FullGuild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(fullGuild.ID)
	if oldGuild != nil {
		oldGuild = &*oldGuild
	}
	guild := disgo.EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyYes)

	eventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Guild:        guild,
		},
		OldGuild: oldGuild,
	})

}
