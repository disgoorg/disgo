package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUserAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUserAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) New() any {
	return &discord.GuildScheduledEventUserEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildScheduledEventUserEvent)

	client.EventManager().Dispatch(&events.GuildScheduledEventUserAddEvent{
		GenericGuildScheduledEventUserEvent: &events.GenericGuildScheduledEventUserEvent{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
