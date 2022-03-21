package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
func (h *gatewayHandlerGuildScheduledEventUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	oldGuildScheduledEvent, _ := client.Caches().GuildScheduledEvents().Get(guildScheduledEvent.GuildID, guildScheduledEvent.ID)
	client.Caches().GuildScheduledEvents().Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)

	client.EventManager().Dispatch(&events.GuildScheduledEventUpdateEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEvent: guildScheduledEvent,
		},
		OldGuildScheduledEvent: oldGuildScheduledEvent,
	})
}
