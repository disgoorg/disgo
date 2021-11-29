package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventUserAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventUserAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) New() interface{} {
	return &discord.GuildScheduledEventUserEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventUserAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildScheduledEventUserEvent)

	bot.EventManager.Dispatch(&events.GuildScheduledEventUserAddEvent{
		GenericGuildScheduledEventUserEvent: &events.GenericGuildScheduledEventUserEvent{
			GenericEvent:          events.NewGenericEvent(bot, sequenceNumber),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
