package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageReactionRemovePayload struct {
	UserID    api.Snowflake  `json:"user_id"`
	ChannelID api.Snowflake  `json:"channel_id"`
	MessageID api.Snowflake  `json:"message_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
	Emoji     *api.Emoji     `json:"emoji"`
}

// MessageReactionRemoveHandler handles api.GatewayEventMessageReactionRemove
type MessageReactionRemoveHandler struct{}

// Event returns the raw gateway event Event
func (h MessageReactionRemoveHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageReactionRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageReactionRemoveHandler) New() interface{} {
	return &messageReactionRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageReactionRemoveHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageReactionRemovePayload)
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

	eventManager.Dispatch(&events.MessageReactionRemoveEvent{
		GenericMessageUserReactionEvent: &events.GenericMessageUserReactionEvent{
			GenericMessageReactionEvent: &events.GenericMessageReactionEvent{
				GenericMessageEvent: genericMessageEvent,
				Emoji:               emoji,
			},
			UserID: payload.UserID,
		},
	})

	if payload.GuildID != nil {
		eventManager.Dispatch(&events.GuildMessageReactionRemoveEvent{
			GenericGuildMessageUserReactionEvent: &events.GenericGuildMessageUserReactionEvent{
				GenericGuildMessageReactionEvent: &events.GenericGuildMessageReactionEvent{
					GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
						GenericMessageEvent: genericMessageEvent,
						GuildID:             *payload.GuildID,
					},
					Emoji: emoji,
				},
				UserID: payload.UserID,
			},
		})

	} else {
		eventManager.Dispatch(&events.DMMessageReactionRemoveEvent{
			GenericDMMessageUserReactionEvent: &events.GenericDMMessageUserReactionEvent{
				GenericDMMessageReactionEvent: &events.GenericDMMessageReactionEvent{
					GenericDMMessageEvent: &events.GenericDMMessageEvent{
						GenericMessageEvent: genericMessageEvent,
					},
					Emoji: emoji,
				},
				UserID: payload.UserID,
			},
		})
	}
}
