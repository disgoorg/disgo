package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerMessageCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.Message)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)
	message := bot.EntityBuilder().CreateMessage(payload, core.CacheStrategyYes)

	bot.EventManager().Dispatch(&events.MessageCreateEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.ID,
			Message:      message,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
		},
	})

	if message.GuildID == nil {
		bot.EventManager().Dispatch(&events.DMMessageCreateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
			},
		})
	} else {
		bot.EventManager().Dispatch(&events.GuildMessageCreateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
			},
		})
	}

}
