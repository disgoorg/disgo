package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

type gatewayHandlerMessageDeleteBulk struct{}

func (h *gatewayHandlerMessageDeleteBulk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDeleteBulk
}

func (h *gatewayHandlerMessageDeleteBulk) New() any {
	return &discord.GatewayEventMessageDeleteBulk{}
}

func (h *gatewayHandlerMessageDeleteBulk) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventMessageDeleteBulk)

	for _, messageID := range payload.IDs {
		handleMessageDelete(client, sequenceNumber, shardID, messageID, payload.ChannelID, payload.GuildID)
	}
}
