package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildUpdateHandler handles api.GuildUpdateGatewayEvent
type GuildUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildUpdateHandler) New() interface{} {
	return discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	guild, ok := v.(discord.Guild)
	if !ok {
		return
	}

	oldCoreGuild := disgo.Caches().GuildCache().GetCopy(guild.ID)

	eventManager.Dispatch(&events.GuildUpdateEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Guild:        disgo.EntityBuilder().CreateGuild(guild, core.CacheStrategyYes),
		},
		OldGuild: oldCoreGuild,
	})

}
