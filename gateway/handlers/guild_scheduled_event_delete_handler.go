package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventDelete) New() any {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	client.Caches().GuildScheduledEvents().Remove(guildScheduledEvent.GuildID, guildScheduledEvent.ID)

	client.EventManager().Dispatch(&events.GuildScheduledEventDeleteEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEvent: guildScheduledEvent,
		},
	})
}
