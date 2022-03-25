package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUserAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUserAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) New() any {
	return &discord.GatewayEventGuildScheduledEventUser{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventGuildScheduledEventUser)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserAddEvent{
		GenericGuildScheduledEventUserEvent: &events.GenericGuildScheduledEventUserEvent{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
