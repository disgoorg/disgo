package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageReactionAddPayload struct {
	UserID    api.Snowflake  `json:"user_id"`
	ChannelID api.Snowflake  `json:"channel_id"`
	MessageID api.Snowflake  `json:"message_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
	Member    *api.Member    `json:"member,omitempty"`
	Emoji     *api.Emoji     `json:"emoji"`
}

// MessageReactionAddHandler handles api.GatewayEventMessageReactionAdd
type MessageReactionAddHandler struct{}

// Event returns the raw gateway event Event
func (h MessageReactionAddHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageReactionAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageReactionAddHandler) New() interface{} {
	return &messageReactionAddPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageReactionAddHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageReactionAddPayload)
	if !ok {
		return
	}

	emoji := disgo.EntityBuilder().CreateEmoji("", payload.Emoji, api.CacheStrategyYes)

	genericMessageEvent := events.GenericMessageEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		Message:      disgo.Cache().Message(payload.ChannelID, payload.MessageID),
	}
	eventManager.Dispatch(genericMessageEvent)

	genericMessageReactionEvent := events.GenericMessageReactionEvent{
		GenericMessageEvent: genericMessageEvent,
		Emoji:               emoji,
	}
	eventManager.Dispatch(genericMessageReactionEvent)

	genericMessageUserReactionEvent := events.GenericMessageUserReactionEvent{
		GenericMessageReactionEvent: genericMessageReactionEvent,
		UserID:                      payload.UserID,
	}
	eventManager.Dispatch(genericMessageUserReactionEvent)

	eventManager.Dispatch(events.MessageReactionAddEvent{
		GenericMessageUserReactionEvent: genericMessageUserReactionEvent,
	})

	if payload.GuildID != nil {
		member := disgo.EntityBuilder().CreateMember(*payload.GuildID, payload.Member, api.CacheStrategyYes)

		genericGuildMessageEvent := events.GenericGuildMessageEvent{
			GenericMessageEvent: genericMessageEvent,
			GuildID:             *payload.GuildID,
		}
		eventManager.Dispatch(genericMessageEvent)

		genericGuildMessageReactionEvent := events.GenericGuildMessageReactionEvent{
			GenericGuildMessageEvent: genericGuildMessageEvent,
			Emoji:                    emoji,
		}
		eventManager.Dispatch(genericGuildMessageReactionEvent)

		genericGuildMessageUserReactionEvent := events.GenericGuildMessageUserReactionEvent{
			GenericGuildMessageReactionEvent: genericGuildMessageReactionEvent,
			UserID:                           payload.UserID,
		}
		eventManager.Dispatch(genericGuildMessageUserReactionEvent)

		eventManager.Dispatch(events.GuildMessageReactionAddEvent{
			GenericGuildMessageUserReactionEvent: genericGuildMessageUserReactionEvent,
			Member:                               member,
		})

	} else {
		genericDMMessageEvent := events.GenericDMMessageEvent{
			GenericMessageEvent: genericMessageEvent,
		}
		eventManager.Dispatch(genericMessageEvent)

		genericDMMessageReactionEvent := events.GenericDMMessageReactionEvent{
			GenericDMMessageEvent: genericDMMessageEvent,
			Emoji:                 emoji,
		}
		eventManager.Dispatch(genericDMMessageReactionEvent)

		genericDMMessageUserReactionEvent := events.GenericDMMessageUserReactionEvent{
			GenericDMMessageReactionEvent: genericDMMessageReactionEvent,
			UserID:                        payload.UserID,
		}
		eventManager.Dispatch(genericDMMessageUserReactionEvent)

		eventManager.Dispatch(events.DMMessageReactionAddEvent{
			GenericDMMessageUserReactionEvent: genericDMMessageUserReactionEvent,
		})
	}
}
