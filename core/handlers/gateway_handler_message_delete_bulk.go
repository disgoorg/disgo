package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageDelete handles discord.GatewayEventTypeMessageDeleteBulk
type gatewayHandlerMessageDeleteBulk struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageDeleteBulk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDeleteBulk
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDeleteBulk) New() interface{} {
	return &discord.MessageDeleteBulkGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDeleteBulk) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.MessageDeleteBulkGatewayEvent)

	for _, messageID := range payload.IDs {
		handleMessageDelete(bot, sequenceNumber, messageID, payload.ChannelID, payload.GuildID)
	}
}
