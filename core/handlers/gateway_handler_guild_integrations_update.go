package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildEmojisUpdate handles discord.GatewayEventTypeGuildIntegrationsUpdate
type gatewayHandlerGuildIntegrationsUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildIntegrationsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildIntegrationsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildIntegrationsUpdate) New() interface{} {
	return &discord.GuildIntegrationsUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildIntegrationsUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GuildIntegrationsUpdateGatewayEvent)

	bot.EventManager().Dispatch(&events.GuildIntegrationsUpdateEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
	})
}
