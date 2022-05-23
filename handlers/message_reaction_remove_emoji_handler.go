package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerMessageReactionRemoveEmoji struct{}

func (h *gatewayHandlerMessageReactionRemoveEmoji) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemoveEmoji
}

func (h *gatewayHandlerMessageReactionRemoveEmoji) New() any {
	return &discord.GatewayEventMessageReactionRemoveEmoji{}
}

func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventMessageReactionRemoveEmoji)

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
