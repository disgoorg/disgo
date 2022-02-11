package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventDelete) New() interface{} {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GuildScheduledEvent)

	bot.Caches.GuildScheduledEvents().Remove(payload.ID)

	bot.EventManager.Dispatch(&events.GuildScheduledEventDeleteEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber),
			GuildScheduledEvent: bot.EntityBuilder.CreateGuildScheduledEvent(payload, core.CacheStrategyNo),
		},
	})
}
