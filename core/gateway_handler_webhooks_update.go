package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerWebhooksUpdate handles core.GatewayEventWebhooksUpdate
type gatewayHandlerWebhooksUpdate struct{}

// EventType returns the raw core.GatewayGatewayEventType
func (h *gatewayHandlerWebhooksUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeWebhooksUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerWebhooksUpdate) New() interface{} {
	return &discord.WebhooksUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerWebhooksUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.WebhooksUpdateGatewayEvent)

	bot.EventManager.Dispatch(&WebhooksUpdateEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		GuildId:      payload.GuildID,
		ChannelID:    payload.ChannelID,
	})
}
