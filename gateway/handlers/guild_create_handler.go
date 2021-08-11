package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildCreateHandler) New() interface{} {
	return &discord.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*discord.FullGuild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(fullGuild.ID)
	wasUnavailable := true
	if oldGuild != nil {
		oldGuild = &*oldGuild
		wasUnavailable = oldGuild.Unavailable
	}
	guild := disgo.EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyYes)

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
