package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageReactionRemoveEmoji struct {}

func (h *gatewayHandlerMessageReactionRemoveEmoji) EventType() gateway.EventType {
	return gateway.EventTypeMessageReactionRemoveEmoji
}

func (h *gatewayHandlerMessageReactionRemoveEmoji) New() any {
	return &gateway.EventMessageReactionRemoveEmoji{}
}

func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventMessageReactionRemoveEmoji)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemoveEmoji{
		GenericEvent: genericEvent,
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		Emoji:        payload.Emoji,
	})

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemoveEmoji{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			Emoji:        payload.Emoji,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemoveEmoji{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      *payload.GuildID,
			Emoji:        payload.Emoji,
		})
	}
}
