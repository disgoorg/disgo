package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemoveAll
type gatewayHandlerMessageReactionRemoveAll struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemoveAll) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemoveAll
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemoveAll) New() any {
	return &discord.GatewayEventMessageReactionRemoveAll{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemoveAll) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	messageReaction := *v.(*discord.GatewayEventMessageReactionRemoveAll)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemoveAll{
		GenericEvent: genericEvent,
		MessageID:    messageReaction.MessageID,
		ChannelID:    messageReaction.ChannelID,
		GuildID:      messageReaction.GuildID,
	})

	if messageReaction.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemoveAll{
			GenericEvent: genericEvent,
			MessageID:    messageReaction.MessageID,
			ChannelID:    messageReaction.ChannelID,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemoveAll{
			GenericEvent: genericEvent,
			MessageID:    messageReaction.MessageID,
			ChannelID:    messageReaction.ChannelID,
			GuildID:      *messageReaction.GuildID,
		})
	}
}
