package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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
	payload := *v.(*discord.Message)

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)
	message := bot.EntityBuilder.CreateMessage(payload, core.CacheStrategyYes)

	bot.EventManager.Dispatch(&events2.MessageCreateEvent{
		GenericMessageEvent: &events2.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.ID,
			Message:      message,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
		},
	})

	if message.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMMessageCreateEvent{
			GenericDMMessageEvent: &events2.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMessageCreateEvent{
			GenericGuildMessageEvent: &events2.GenericGuildMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
			},
		})
	}

}
