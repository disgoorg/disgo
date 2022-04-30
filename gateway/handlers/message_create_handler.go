package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerMessageCreate handles discord.GatewayEventTypeMessageCreate
type gatewayHandlerMessageCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageCreate) New() any {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	message := *v.(*discord.Message)

	if message.Member != nil {
		message.Member.User = message.Author
	}

	client.Caches().Messages().Put(message.ChannelID, message.ID, message)

	// update discord.MessageChannel.LastMessageID()
	if channel, ok := client.Caches().Channels().GetMessageChannel(message.ChannelID); ok {
		client.Caches().Channels().Put(message.ChannelID, discord.ApplyLastMessageID(channel, message.ID))
	}

	genericEvent := events.NewGenericEvent(client, sequenceNumber)
	client.EventManager().DispatchEvent(&events.MessageCreateEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    message.ID,
			Message:      message,
			ChannelID:    message.ChannelID,
			GuildID:      message.GuildID,
		},
	})

	if message.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageCreateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageCreateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
				GuildID:      *message.GuildID,
			},
		})
	}

}
