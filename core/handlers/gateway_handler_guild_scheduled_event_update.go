package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUpdate) New() interface{} {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.GuildScheduledEvent)

	oldGuildScheduledEvent := bot.Caches.GuildScheduledEvents().GetCopy(payload.ID)

	bot.EventManager.Dispatch(&events.GuildScheduledEventUpdateEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber, shardID),
			GuildScheduledEvent: bot.EntityBuilder.CreateGuildScheduledEvent(payload, core.CacheStrategyYes),
		},
		OldGuildScheduledEvent: oldGuildScheduledEvent,
	})
}
