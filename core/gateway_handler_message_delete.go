package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type messageDeletePayload struct {
	MessageID discord.Snowflake  `json:"id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	ChannelID discord.Snowflake  `json:"channel_id"`
}

// gatewayHandlerMessageDelete handles core.GatewayEventMessageDelete
type gatewayHandlerMessageDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerMessageDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageDelete) New() interface{} {
	return &messageDeletePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*messageDeletePayload)

	genericMessageEvent := &GenericMessageEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MessageID:    payload.MessageID,
		Message:      bot.Caches.MessageCache().GetCopy(payload.ChannelID, payload.MessageID),
		ChannelID:    payload.ChannelID,
	}

	bot.Caches.MessageCache().Remove(payload.ChannelID, payload.MessageID)

	bot.EventManager.Dispatch(&MessageDeleteEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&DMMessageDeleteEvent{
			GenericDMMessageEvent: &GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&GuildMessageDeleteEvent{
			GenericGuildMessageEvent: &GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *payload.GuildID,
			},
		})
	}
}
