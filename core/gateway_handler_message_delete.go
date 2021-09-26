package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageDelete handles core.GatewayEventMessageDelete
type gatewayHandlerMessageDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerMessageDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDelete) New() interface{} {
	return &discord.MessageDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.MessageDeleteGatewayEvent)

	handleMessageDelete(bot, sequenceNumber, payload.ID, payload.ChannelID, payload.GuildID)
}

func handleMessageDelete(bot *Bot, sequenceNumber int, messageID discord.Snowflake, channelID discord.Snowflake, guildID *discord.Snowflake) {
	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MessageID:    messageID,
		Message:      bot.Caches.MessageCache().GetCopy(channelID, messageID),
		ChannelID:    channelID,
	}

	bot.Caches.MessageCache().Remove(channelID, messageID)

	bot.EventManager.Dispatch(&MessageDeleteEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if guildID == nil {
		bot.EventManager.Dispatch(&DMMessageDeleteEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageDeleteEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *guildID,
			},
		})
	}
}
