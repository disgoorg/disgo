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
	return &discord.GatewayEventMessageDeleteBulk{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDeleteBulk) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventMessageDeleteBulk)

	for _, messageID := range payload.IDs {
		handleMessageDelete(client, sequenceNumber, shardID, messageID, payload.ChannelID, payload.GuildID)
	}
}
