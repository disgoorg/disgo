package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemoveEmoji
type gatewayHandlerMessageReactionRemoveEmoji struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemoveEmoji) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemoveEmoji
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemoveEmoji) New() interface{} {
	return &discord.GatewayEventMessageReactionRemoveEmoji{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventMessageReactionRemoveEmoji)

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events2.MessageReactionRemoveEmojiEvent{
		GenericEvent: genericEvent,
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		Emoji:        payload.Emoji,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMMessageReactionRemoveEmojiEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			Emoji:        payload.Emoji,
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMessageReactionRemoveEmojiEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      *payload.GuildID,
			Emoji:        payload.Emoji,
		})
	}
}
