package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemoveEmoji
type gatewayHandlerMessageReactionRemoveEmoji struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemoveEmoji) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemoveEmoji
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemoveEmoji) New() any {
	return &discord.GatewayEventMessageReactionRemoveEmoji{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventMessageReactionRemoveEmoji)

	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	client.EventManager().Dispatch(&events.MessageReactionRemoveEmojiEvent{
		GenericEvent: genericEvent,
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		Emoji:        payload.Emoji,
	})

	if payload.GuildID == nil {
		client.EventManager().Dispatch(&events.DMMessageReactionRemoveEmojiEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			Emoji:        payload.Emoji,
		})
	} else {
		client.EventManager().Dispatch(&events.GuildMessageReactionRemoveEmojiEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      *payload.GuildID,
			Emoji:        payload.Emoji,
		})
	}
}
