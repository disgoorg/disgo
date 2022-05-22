package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// gatewayHandlerMessageDelete handles discord.GatewayEventTypeMessageDelete
type gatewayHandlerMessageDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDelete) New() any {
	return &discord.GatewayEventMessageDelete{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventMessageDelete)

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
