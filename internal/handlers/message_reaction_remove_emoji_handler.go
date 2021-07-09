package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageReactionRemoveEmotePayload struct {
	ChannelID api.Snowflake  `json:"channel_id"`
	MessageID api.Snowflake  `json:"message_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
	Emoji     *api.Emoji     `json:"emoji"`
}

// MessageReactionRemoveEmoteHandler handles api.GatewayEventMessageReactionRemoveEmoji
type MessageReactionRemoveEmoteHandler struct{}

// Event returns the raw gateway event Event
func (h *MessageReactionRemoveEmoteHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageReactionRemoveEmoji
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageReactionRemoveEmoteHandler) New() interface{} {
	return &messageReactionRemoveEmotePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageReactionRemoveEmoteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageReactionRemoveEmotePayload)
	if !ok {
		return
	}

	emoji := disgo.EntityBuilder().CreateEmoji("", payload.Emoji, api.CacheStrategyYes)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		Message:      disgo.Cache().Message(payload.ChannelID, payload.MessageID),
	}

	eventManager.Dispatch(&events.MessageReactionRemoveEmojiEvent{
		GenericMessageReactionEvent: &events.GenericMessageReactionEvent{
			GenericMessageEvent: genericMessageEvent,
			Emoji:               emoji,
		},
	})

	if payload.GuildID != nil {
		eventManager.Dispatch(&events.GuildMessageReactionRemoveEmojiEvent{
			GenericGuildMessageReactionEvent: &events.GenericGuildMessageReactionEvent{
				GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
					GenericMessageEvent: genericMessageEvent,
					GuildID:             *payload.GuildID,
				},
				Emoji: emoji,
			},
		})

	} else {
		eventManager.Dispatch(&events.DMMessageReactionRemoveEmojiEvent{
			GenericDMMessageReactionEvent: &events.GenericDMMessageReactionEvent{
				GenericDMMessageEvent: &events.GenericDMMessageEvent{
					GenericMessageEvent: genericMessageEvent,
				},
				Emoji: emoji,
			},
		})
	}
}
