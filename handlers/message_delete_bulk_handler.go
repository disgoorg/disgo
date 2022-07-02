package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageDeleteBulk struct {}

func (h *gatewayHandlerMessageDeleteBulk) EventType() gateway.EventType {
	return gateway.EventTypeMessageDeleteBulk
}

func (h *gatewayHandlerMessageDeleteBulk) New() any {
	return &gateway.EventMessageDeleteBulk{}
}

func (h *gatewayHandlerMessageDeleteBulk) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventMessageDeleteBulk)

	for _, messageID := range payload.IDs {
		handleMessageDelete(client, sequenceNumber, shardID, messageID, payload.ChannelID, payload.GuildID)
	}
}
