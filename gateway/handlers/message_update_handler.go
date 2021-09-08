package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
	return discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	message, ok := v.(discord.Message)
	if !ok {
		return
	}

	oldCoreMessage := bot.Caches.MessageCache().GetCopy(message.ChannelID, message.ID)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		Message:      bot.EntityBuilder.CreateMessage(message, core.CacheStrategyYes),
	}

	bot.EventManager.Dispatch(&events.MessageUpdateEvent{
		GenericMessageEvent: genericMessageEvent,
		OldMessage:          oldCoreMessage,
	})

	if message.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageUpdateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			OldMessage: oldCoreMessage,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
			OldMessage: oldCoreMessage,
		})
	}
}
