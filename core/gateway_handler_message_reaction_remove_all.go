package core

import "github.com/DisgoOrg/disgo/discord"

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemoveAll
type gatewayHandlerMessageReactionRemoveAll struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemoveAll) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemoveAll
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemoveAll) New() interface{} {
	return &discord.GatewayEventMessageReactionRemoveAll{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemoveAll) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemoveAll)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageReaction.MessageID,
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
		ChannelID:    messageReaction.ChannelID,
	}
	bot.EventManager.Dispatch(&MessageReactionRemoveAllEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageReactionRemoveAllEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageReactionRemoveAllEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *messageReaction.GuildID,
			},
		})
	}
}
