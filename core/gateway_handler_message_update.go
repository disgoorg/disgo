package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// MessageUpdateHandler handles api.GatewayEventMessageUpdate
type MessageUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *MessageUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageUpdateHandler) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageUpdateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	message := *v.(*discord.Message)

	oldCoreMessage := bot.Caches.MessageCache().GetCopy(message.ChannelID, message.ID)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		Message:      bot.EntityBuilder.CreateMessage(message, CacheStrategyYes),
	}

	bot.EventManager.Dispatch(&MessageUpdateEvent{
		GenericMessageEvent: genericMessageEvent,
		OldMessage:          oldCoreMessage,
	})

	if message.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageUpdateEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			OldMessage: oldCoreMessage,
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
			OldMessage: oldCoreMessage,
		})
	}
}
