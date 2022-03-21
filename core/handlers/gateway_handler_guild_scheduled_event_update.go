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
func (h *gatewayHandlerGuildScheduledEventUpdate) New() any {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	oldGuildScheduledEvent, _ := bot.Caches().GuildScheduledEvents().Get(guildScheduledEvent.GuildID, guildScheduledEvent.ID)
	bot.Caches().GuildScheduledEvents().Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)

	bot.EventManager().Dispatch(&events.GuildScheduledEventUpdateEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber),
			GuildScheduledEvent: guildScheduledEvent,
		},
		OldGuildScheduledEvent: oldGuildScheduledEvent,
	})
}
