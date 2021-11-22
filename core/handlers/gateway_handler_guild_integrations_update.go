package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerGuildIntegrationsUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildIntegrationsUpdateGatewayEvent)

	bot.EventManager.Dispatch(&events2.GuildIntegrationsUpdateEvent{
		GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
	})
}
