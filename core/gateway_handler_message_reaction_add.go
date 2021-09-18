package core

import "github.com/DisgoOrg/disgo/discord"

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionAdd
type gatewayHandlerMessageReactionAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionAdd) New() interface{} {
	return &discord.GatewayEventMessageReactionAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionAdd) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionAdd)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
	}
	bot.EventManager.Dispatch(&MessageReactionAddEvent{
		GenericReactionEvent: &GenericReactionEvent{
			GenericMessageEvent: genericMessageEvent,
			UserID:              messageReaction.UserID,
			Emoji:               messageReaction.Emoji,
		},
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageReactionAddEvent{
			GenericDMMessageReactionEvent: &GenericDMMessageReactionEvent{
				GenericDMMessageEvent: &GenericDMMessageEvent{
					GenericMessageEvent: genericMessageEvent,
				},
				UserID: messageReaction.UserID,
				Emoji:  messageReaction.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageReactionAddEvent{
			GenericGuildMessageReactionEvent: &GenericGuildMessageReactionEvent{
				GenericGuildMessageEvent: &GenericGuildMessageEvent{
					GenericMessageEvent: genericMessageEvent,
					GuildID:             *messageReaction.GuildID,
				},
				UserID: messageReaction.UserID,
				Member: bot.EntityBuilder.CreateMember(*messageReaction.GuildID, *messageReaction.Member, CacheStrategyYes),
				Emoji:  messageReaction.Emoji,
			},
		})
	}
}
