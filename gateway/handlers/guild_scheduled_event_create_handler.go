package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventCreate) New() any {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventCreate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	client.Caches().GuildScheduledEvents().Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)

	client.EventManager().Dispatch(&events.GuildScheduledEventCreateEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEvent: guildScheduledEvent,
		},
	})
}
