package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *GuildUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildUpdateHandler) New() interface{} {
	return discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	guild, ok := i.(discord.Guild)
	if !ok {
		return
	}

	oldCoreGuild := disgo.Cache().GuildCache().GetCopy(guild.ID)

	eventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Guild:        disgo.EntityBuilder().CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldCoreGuild,
	})

}
