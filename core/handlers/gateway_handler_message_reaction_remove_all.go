package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events2.MessageReactionRemoveAllEvent{
		GenericEvent: genericEvent,
		MessageID:    messageReaction.MessageID,
		ChannelID:    messageReaction.ChannelID,
		GuildID:      messageReaction.GuildID,
	})

	if messageReaction.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMMessageReactionRemoveAllEvent{
			GenericEvent: genericEvent,
			MessageID:    messageReaction.MessageID,
			ChannelID:    messageReaction.ChannelID,
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMessageReactionRemoveAllEvent{
			GenericEvent: genericEvent,
			MessageID:    messageReaction.MessageID,
			ChannelID:    messageReaction.ChannelID,
			GuildID:      *messageReaction.GuildID,
		})
	}
}
