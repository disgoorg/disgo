package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageUpdate
type gatewayHandlerMessageUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageUpdate) New() any {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	message := *v.(*discord.Message)

	oldMessage, _ := client.Caches().Messages().Get(message.ChannelID, message.ID)
	client.Caches().Messages().Put(message.ChannelID, message.ID, message)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	client.EventManager().DispatchEvent(&events.MessageUpdate{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    message.ID,
			Message:      message,
			ChannelID:    message.ChannelID,
			GuildID:      message.GuildID,
		},
		OldMessage: oldMessage,
	})

	if message.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageUpdate{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
			},
			OldMessage: oldMessage,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageUpdate{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
				GuildID:      *message.GuildID,
			},
			OldMessage: oldMessage,
		})
	}
}
