package core

import "github.com/DisgoOrg/disgo/discord"

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
func (h *gatewayHandlerMessageReactionRemoveEmoji) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemoveEmoji)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
	}
	bot.EventManager.Dispatch(&MessageReactionRemoveEmojiEvent{
		GenericMessageEvent: genericMessageEvent,
		Emoji:               messageReaction.Emoji,
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageReactionRemoveEmojiEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			Emoji: messageReaction.Emoji,
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageReactionRemoveEmojiEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *messageReaction.GuildID,
			},
			Emoji: messageReaction.Emoji,
		})
	}
}
