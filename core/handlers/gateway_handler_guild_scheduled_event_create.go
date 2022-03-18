package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildScheduledEventCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildScheduledEventCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildScheduledEventCreate) New() interface{} {
	return &discord.GuildScheduledEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildScheduledEventCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GuildScheduledEvent)

	bot.EventManager().Dispatch(&events.GuildScheduledEventCreateEvent{
		GenericGuildScheduledEventEvent: &events.GenericGuildScheduledEventEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber),
			GuildScheduledEvent: bot.EntityBuilder.CreateGuildScheduledEvent(payload, core.CacheStrategyYes),
		},
	})
}
