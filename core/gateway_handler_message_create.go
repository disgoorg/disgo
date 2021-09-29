package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageCreate handles discord.GatewayEventTypeMessageCreate
type gatewayHandlerMessageCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageCreate) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageCreate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	message := *v.(*discord.Message)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MessageID:    message.ID,
		Message:      bot.EntityBuilder.CreateMessage(message, CacheStrategyYes),
		ChannelID:    message.ChannelID,
	}

	bot.EventManager.Dispatch(&MessageCreateEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if message.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageCreateEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageCreateEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
		})
	}

}
