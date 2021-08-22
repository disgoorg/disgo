package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildCreateHandler) New() interface{} {
	return discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	guild, ok := i.(discord.Guild)
	if !ok {
		return
	}

	oldCoreGuild := disgo.Cache().GuildCache().GetCopy(guild.ID)
	wasUnavailable := true
	if oldCoreGuild != nil {
		wasUnavailable = oldCoreGuild.Unavailable
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		GuildID:      guild.ID,
		Guild:        disgo.EntityBuilder().CreateGuild(guild, core.CacheStrategyYes),
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
