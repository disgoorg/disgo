package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// Event returns the raw gateway event Event
func (h *GuildCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildCreateHandler) New() interface{} {
	return &api.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*api.FullGuild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(fullGuild.ID)
	wasUnavailable := true
	if oldGuild != nil {
		oldGuild = &*oldGuild
		wasUnavailable = oldGuild.Unavailable
	}
	guild := disgo.EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyYes)

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		GuildID:      guild.ID,
		Guild:        guild,
	}

	if !guild.Ready {
		guild.Ready = true
		eventManager.Dispatch(&events.GuildReadyEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}

	if wasUnavailable {
		eventManager.Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		eventManager.Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
