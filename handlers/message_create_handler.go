package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerMessageCreate struct{}

func (h *gatewayHandlerMessageCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageCreate
}

func (h *gatewayHandlerMessageCreate) New() any {
	return &discord.Message{}
}

func (h *gatewayHandlerMessageCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	message := *v.(*discord.Message)

	if message.Flags.Has(discord.MessageFlagEphemeral) {
		// Ignore ephemeral messages as they miss guild_id & member
		return
	}

	if message.Member != nil {
		message.Member.User = message.Author
	}

	client.Caches().Messages().Put(message.ChannelID, message.ID, message)

	if channel, ok := client.Caches().Channels().GetMessageChannel(message.ChannelID); ok {
		client.Caches().Channels().Put(message.ChannelID, discord.ApplyLastMessageIDToChannel(channel, message.ID))
	}

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	client.EventManager().DispatchEvent(&events.MessageCreate{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    message.ID,
			Message:      message,
			ChannelID:    message.ChannelID,
			GuildID:      message.GuildID,
		},
	})

	if message.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageCreate{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageCreate{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    message.ID,
				Message:      message,
				ChannelID:    message.ChannelID,
				GuildID:      *message.GuildID,
			},
		})
	}

}
