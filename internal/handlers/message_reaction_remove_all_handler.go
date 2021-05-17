package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageReactionRemoveAllPayload struct {
	ChannelID api.Snowflake  `json:"channel_id"`
	MessageID api.Snowflake  `json:"message_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
}

// MessageReactionRemoveAllHandler handles api.GatewayEventMessageReactionRemove
type MessageReactionRemoveAllHandler struct{}

// Event returns the raw gateway event Event
func (h MessageReactionRemoveAllHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageReactionRemoveAll
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageReactionRemoveAllHandler) New() interface{} {
	return &messageReactionRemoveAllPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageReactionRemoveAllHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageReactionRemoveAllPayload)
	if !ok {
		return
	}

	genericMessageEvent := events.GenericMessageEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		Message:      disgo.Cache().Message(payload.ChannelID, payload.MessageID),
	}
	eventManager.Dispatch(genericMessageEvent)

	eventManager.Dispatch(events.MessageReactionRemoveAllEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if payload.GuildID != nil {
		genericGuildMessageEvent := events.GenericGuildMessageEvent{
			GenericMessageEvent: genericMessageEvent,
			GuildID:             *payload.GuildID,
		}
		eventManager.Dispatch(genericMessageEvent)

		eventManager.Dispatch(events.GuildMessageReactionRemoveAllEvent{
			GenericGuildMessageEvent: genericGuildMessageEvent,
		})

	} else {
		genericDMMessageEvent := events.GenericDMMessageEvent{
			GenericMessageEvent: genericMessageEvent,
		}
		eventManager.Dispatch(genericDMMessageEvent)

		eventManager.Dispatch(events.DMMessageReactionRemoveAllEvent{
			GenericDMMessageEvent: genericDMMessageEvent,
		})
	}
}
