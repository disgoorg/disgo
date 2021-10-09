package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
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
func (h *gatewayHandlerMessageReactionRemoveEmoji) New() interface{} {
	return &discord.GatewayEventMessageReactionRemoveEmoji{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemoveEmoji)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageReaction.MessageID,
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
		ChannelID:    messageReaction.ChannelID,
	}
	bot.EventManager.Dispatch(&events.MessageReactionRemoveEmojiEvent{
		GenericMessageEvent: genericMessageEvent,
		Emoji:               messageReaction.Emoji,
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageReactionRemoveEmojiEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			Emoji: messageReaction.Emoji,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageReactionRemoveEmojiEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *messageReaction.GuildID,
			},
			Emoji: messageReaction.Emoji,
		})
	}
}
