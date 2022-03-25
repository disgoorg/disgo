package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUserRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUserRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserRemove) New() any {
	return &discord.GatewayEventGuildScheduledEventUserEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventGuildScheduledEventUserEvent)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserRemoveEvent{
		GenericGuildScheduledEventUserEvent: &events.GenericGuildScheduledEventUserEvent{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
