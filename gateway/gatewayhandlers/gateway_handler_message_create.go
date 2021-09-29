package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerMessageCreate handles core.GatewayEventMessageCreate
type gatewayHandlerMessageCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerMessageCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageCreate) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	message := *v.(*discord.Message)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		MessageID:    message.ID,
		Message:      bot.EntityBuilder.CreateMessage(message, core.CacheStrategyYes),
		ChannelID:    message.ChannelID,
	}

	bot.EventManager.Dispatch(&events.MessageCreateEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if message.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageCreateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageCreateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
		})
	}

}
