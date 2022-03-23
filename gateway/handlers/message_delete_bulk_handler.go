package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

// gatewayHandlerMessageDelete handles discord.GatewayEventTypeMessageDeleteBulk
type gatewayHandlerMessageDeleteBulk struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageDeleteBulk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDeleteBulk
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDeleteBulk) New() any {
	return &discord.MessageDeleteBulkGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDeleteBulk) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.MessageDeleteBulkGatewayEvent)

	for _, messageID := range payload.IDs {
		handleMessageDelete(client, sequenceNumber, messageID, payload.ChannelID, payload.GuildID)
	}
}
