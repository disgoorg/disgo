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
	Emote     *api.Emote     `json:"emoji"`
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

	emote := disgo.EntityBuilder().CreateEmote("", payload.Emote, api.CacheStrategyYes)

	genericMessageEvent := events.GenericMessageEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		Message:      disgo.Cache().Message(payload.ChannelID, payload.MessageID),
	}
	eventManager.Dispatch(genericMessageEvent)

	genericMessageReactionEvent := events.GenericMessageReactionEvent{
		GenericMessageEvent: genericMessageEvent,
		Emote:               emote,
	}
	eventManager.Dispatch(genericMessageReactionEvent)

	genericMessageUserReactionEvent := events.GenericMessageUserReactionEvent{
		GenericMessageReactionEvent: genericMessageReactionEvent,
		UserID:                      payload.UserID,
	}
	eventManager.Dispatch(genericMessageUserReactionEvent)

	eventManager.Dispatch(events.MessageReactionRemoveEvent{
		GenericMessageUserReactionEvent: genericMessageUserReactionEvent,
	})

	if payload.GuildID != nil {
		genericGuildMessageEvent := events.GenericGuildMessageEvent{
			GenericMessageEvent: genericMessageEvent,
			GuildID:             *payload.GuildID,
		}
		eventManager.Dispatch(genericMessageEvent)

		genericGuildMessageReactionEvent := events.GenericGuildMessageReactionEvent{
			GenericGuildMessageEvent: genericGuildMessageEvent,
			Emote:                    emote,
		}
		eventManager.Dispatch(genericGuildMessageReactionEvent)

		genericGuildMessageUserReactionEvent := events.GenericGuildMessageUserReactionEvent{
			GenericGuildMessageReactionEvent: genericGuildMessageReactionEvent,
			UserID:                           payload.UserID,
		}
		eventManager.Dispatch(genericGuildMessageUserReactionEvent)

		eventManager.Dispatch(events.GuildMessageReactionRemoveEvent{
			GenericGuildMessageUserReactionEvent: genericGuildMessageUserReactionEvent,
		})

	} else {
		genericDMMessageEvent := events.GenericDMMessageEvent{
			GenericMessageEvent: genericMessageEvent,
		}
		eventManager.Dispatch(genericMessageEvent)

		genericDMMessageReactionEvent := events.GenericDMMessageReactionEvent{
			GenericDMMessageEvent: genericDMMessageEvent,
			Emote:                 emote,
		}
		eventManager.Dispatch(genericDMMessageReactionEvent)

		genericDMMessageUserReactionEvent := events.GenericDMMessageUserReactionEvent{
			GenericDMMessageReactionEvent: genericDMMessageReactionEvent,
			UserID:                        payload.UserID,
		}
		eventManager.Dispatch(genericDMMessageUserReactionEvent)

		eventManager.Dispatch(events.DMMessageReactionRemoveEvent{
			GenericDMMessageUserReactionEvent: genericDMMessageUserReactionEvent,
		})
	}
}
