package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events.MessageUpdateEvent{
		GenericMessageEvent: &events.GenericMessageEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.ID,
			Message:      message,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
		},
		OldMessage: oldMessage,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageUpdateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.ID,
				Message:      message,
				ChannelID:    payload.ChannelID,
			},
			OldMessage: oldMessage,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
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
