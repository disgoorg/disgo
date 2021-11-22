package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageUpdate handles core.GatewayEventMessageUpdate
type gatewayHandlerMessageUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerMessageUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageUpdate) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.Message)

	oldMessage := bot.Caches.Messages().GetCopy(payload.ChannelID, payload.ID)

	message := bot.EntityBuilder.CreateMessage(payload, core.CacheStrategyYes)

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events2.MessageUpdateEvent{
		GenericMessageEvent: &events2.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.ID,
			Message:      message,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
		},
		OldMessage: oldMessage,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMMessageUpdateEvent{
			GenericDMMessageEvent: &events2.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
			},
			OldMessage: oldMessage,
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &events2.GenericGuildMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
			},
			OldMessage: oldMessage,
		})
	}
}
