package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

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
func (h *gatewayHandlerMessageReactionAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionAdd)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageReaction.MessageID,
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
		ChannelID:    messageReaction.ChannelID,
	}
	bot.EventManager.Dispatch(&events.MessageReactionAddEvent{
		GenericReactionEvent: &events.GenericReactionEvent{
			GenericMessageEvent: genericMessageEvent,
			UserID:              messageReaction.UserID,
			Emoji:               messageReaction.Emoji,
		},
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageReactionAddEvent{
			GenericDMMessageReactionEvent: &events.GenericDMMessageReactionEvent{
				GenericDMMessageEvent: &events.GenericDMMessageEvent{
					GenericMessageEvent: genericMessageEvent,
				},
				UserID: messageReaction.UserID,
				Emoji:  messageReaction.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageReactionAddEvent{
			GenericGuildMessageReactionEvent: &events.GenericGuildMessageReactionEvent{
				GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
					GenericMessageEvent: genericMessageEvent,
					GuildID:             *messageReaction.GuildID,
				},
				UserID: messageReaction.UserID,
				Member: bot.EntityBuilder.CreateMember(*messageReaction.GuildID, *messageReaction.Member, core.CacheStrategyYes),
				Emoji:  messageReaction.Emoji,
			},
		})
	}
}
