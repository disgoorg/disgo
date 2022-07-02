package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageReactionRemove struct {}

func (h *gatewayHandlerMessageReactionRemove) EventType() gateway.EventType {
	return gateway.EventTypeMessageReactionRemove
}

func (h *gatewayHandlerMessageReactionRemove) New() any {
	return &gateway.EventMessageReactionRemove{}
}

func (h *gatewayHandlerMessageReactionRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventMessageReactionRemove)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemove{
		GenericReaction: &events.GenericReaction{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
			UserID:       payload.UserID,
			Emoji:        payload.Emoji,
		},
	})

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemove{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemove{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	}
}
