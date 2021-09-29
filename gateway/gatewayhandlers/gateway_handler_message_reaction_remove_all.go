package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

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
func (h *gatewayHandlerMessageReactionRemoveAll) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemoveAll)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageReaction.MessageID,
		Message:      bot.Caches.MessageCache().Get(messageReaction.ChannelID, messageReaction.MessageID),
		ChannelID:    messageReaction.ChannelID,
	}
	bot.EventManager.Dispatch(&events.MessageReactionRemoveAllEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageReactionRemoveAllEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageReactionRemoveAllEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *messageReaction.GuildID,
			},
		})
	}
}
