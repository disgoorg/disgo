package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

func gatewayHandlerMessageDelete struct {}

func (h *gatewayHandlerMessageDelete) EventType() gateway.EventType {
	return gateway.EventTypeMessageDelete
}

func (h *gatewayHandlerMessageDelete) New() any {
	return &gateway.EventMessageDelete{}
}

func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventMessageDelete)

	handleMessageDelete(client, sequenceNumber, shardID, payload.ID, payload.ChannelID, payload.GuildID)
}

func handleMessageDelete(client bot.Client, sequenceNumber int, shardID int, messageID snowflake.ID, channelID snowflake.ID, guildID *snowflake.ID) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	message, _ := client.Caches().Messages().Remove(channelID, messageID)

	client.EventManager().DispatchEvent(&events.MessageDelete{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    messageID,
			Message:      message,
			ChannelID:    channelID,
		},
	})

	if guildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageDelete{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageDelete{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
				GuildID:      *guildID,
			},
		})
	}
}
