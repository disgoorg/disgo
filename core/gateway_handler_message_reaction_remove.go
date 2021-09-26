package core

import "github.com/DisgoOrg/disgo/discord"

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemove
type gatewayHandlerMessageReactionRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemove) New() interface{} {
	return &discord.GatewayEventMessageReactionRemove{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemove) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemove)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageReaction.MessageID,
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
		ChannelID:    messageReaction.ChannelID,
	}
	bot.EventManager.Dispatch(&MessageReactionRemoveEvent{
		GenericReactionEvent: &GenericReactionEvent{
			GenericMessageEvent: genericMessageEvent,
			UserID:              messageReaction.UserID,
			Emoji:               messageReaction.Emoji,
		},
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageReactionRemoveEvent{
			GenericDMMessageReactionEvent: &GenericDMMessageReactionEvent{
				GenericDMMessageEvent: &GenericDMMessageEvent{
					GenericMessageEvent: genericMessageEvent,
				},
				UserID: messageReaction.UserID,
				Emoji:  messageReaction.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageReactionRemoveEvent{
			GenericGuildMessageReactionEvent: &GenericGuildMessageReactionEvent{
				GenericGuildMessageEvent: &GenericGuildMessageEvent{
					GenericMessageEvent: genericMessageEvent,
					GuildID:             *messageReaction.GuildID,
				},
				UserID: messageReaction.UserID,
				Member: bot.Caches.MemberCache().Get(*messageReaction.GuildID, messageReaction.UserID),
				Emoji:  messageReaction.Emoji,
			},
		})
	}
}
