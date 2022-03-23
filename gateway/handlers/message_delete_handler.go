package handlers

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerMessageDelete handles discord.GatewayEventTypeMessageDelete
type gatewayHandlerMessageDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDelete) New() any {
	return &discord.MessageDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.MessageDeleteGatewayEvent)

	handleMessageDelete(client, sequenceNumber, payload.ID, payload.ChannelID, payload.GuildID)
}

func handleMessageDelete(client bot.Client, sequenceNumber discord.GatewaySequence, messageID snowflake.Snowflake, channelID snowflake.Snowflake, guildID *snowflake.Snowflake) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	message, _ := client.Caches().Messages().Remove(channelID, messageID)

	client.EventManager().Dispatch(&events.MessageDeleteEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    messageID,
			Message:      message,
			ChannelID:    channelID,
		},
	})

	if guildID == nil {
		client.EventManager().Dispatch(&events.DMMessageDeleteEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
			},
		})
	} else {
		client.EventManager().Dispatch(&events.GuildMessageDeleteEvent{
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
