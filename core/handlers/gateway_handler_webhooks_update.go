package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerWebhooksUpdate handles discord.GatewayEventTypeWebhooksUpdate
type gatewayHandlerWebhooksUpdate struct{}

// EventType returns the raw discord.GatewayEventType
func (h *gatewayHandlerWebhooksUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeWebhooksUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerWebhooksUpdate) New() interface{} {
	return &discord.WebhooksUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerWebhooksUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.WebhooksUpdateGatewayEvent)

	bot.EventManager.Dispatch(&events.WebhooksUpdateEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
		GuildId:      payload.GuildID,
		ChannelID:    payload.ChannelID,
	})
}
