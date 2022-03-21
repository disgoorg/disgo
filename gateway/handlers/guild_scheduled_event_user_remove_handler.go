package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUserRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUserRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserRemove) New() any {
	return &discord.GuildScheduledEventUserEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserRemove) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildScheduledEventUserEvent)

	client.EventManager().Dispatch(&events.GuildScheduledEventUserRemoveEvent{
		GenericGuildScheduledEventUserEvent: &events.GenericGuildScheduledEventUserEvent{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
