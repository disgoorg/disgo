package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageReactionRemoveEmotePayload struct {
	ChannelID api.Snowflake  `json:"channel_id"`
	MessageID api.Snowflake  `json:"message_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
	Emote     *api.Emote     `json:"emoji"`
}

// MessageReactionRemoveEmoteHandler handles api.GatewayEventMessageReactionRemoveEmoji
type MessageReactionRemoveEmoteHandler struct{}

// Event returns the raw gateway event Event
func (h MessageReactionRemoveEmoteHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageReactionRemoveEmoji
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageReactionRemoveEmoteHandler) New() interface{} {
	return &messageReactionRemoveEmotePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageReactionRemoveEmoteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageReactionRemoveEmotePayload)
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

	eventManager.Dispatch(events.MessageReactionRemoveEmoteEvent{
		GenericMessageReactionEvent: genericMessageReactionEvent,
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

		eventManager.Dispatch(events.GuildMessageReactionRemoveEmoteEvent{
			GenericGuildMessageReactionEvent: genericGuildMessageReactionEvent,
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

		eventManager.Dispatch(events.DMMessageReactionRemoveEmoteEvent{
			GenericDMMessageReactionEvent: genericDMMessageReactionEvent,
		})
	}
}
