package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// gatewayHandlerMessageDelete handles discord.GatewayEventTypeMessageDelete
type gatewayHandlerMessageDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDelete) New() interface{} {
	return &discord.MessageDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.MessageDeleteGatewayEvent)

	handleMessageDelete(bot, sequenceNumber, shardID, payload.ID, payload.ChannelID, payload.GuildID)
}

func handleMessageDelete(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, messageID snowflake.Snowflake, channelID snowflake.Snowflake, guildID *snowflake.Snowflake) {
	genericEvent := events.NewGenericEvent(bot, sequenceNumber, shardID)

	message := bot.Caches.Messages().GetCopy(channelID, messageID)
	bot.Caches.Messages().Remove(channelID, messageID)

	bot.EventManager.Dispatch(&events.MessageDeleteEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    messageID,
			Message:      message,
			ChannelID:    channelID,
		},
	})

	if guildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageDeleteEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageDeleteEvent{
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
