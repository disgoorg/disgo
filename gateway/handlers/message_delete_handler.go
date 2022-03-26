package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake"
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
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventMessageDelete)

	handleMessageDelete(client, sequenceNumber, payload.ID, payload.ChannelID, payload.GuildID)
}

func handleMessageDelete(client bot.Client, sequenceNumber int, messageID snowflake.Snowflake, channelID snowflake.Snowflake, guildID *snowflake.Snowflake) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	message, _ := client.Caches().Messages().Remove(channelID, messageID)

	client.EventManager().DispatchEvent(&events.MessageDeleteEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    messageID,
			Message:      message,
			ChannelID:    channelID,
		},
	})

	if guildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageDeleteEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageDeleteEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
				GuildID:      *guildID,
			},
		})
	}
}
