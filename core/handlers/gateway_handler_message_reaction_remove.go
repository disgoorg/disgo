package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionRemove
type gatewayHandlerMessageReactionRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionRemove) New() interface{} {
	return &discord.GatewayEventMessageReactionRemove{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionRemove) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventMessageReactionRemove)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events.MessageReactionRemoveEvent{
		GenericReactionEvent: &events.GenericReactionEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
			UserID:       payload.UserID,
			Emoji:        payload.Emoji,
		},
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMMessageReactionRemoveEvent{
			GenericDMMessageReactionEvent: &events.GenericDMMessageReactionEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageReactionRemoveEvent{
			GenericGuildMessageReactionEvent: &events.GenericGuildMessageReactionEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	}
}
