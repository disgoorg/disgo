package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageDeletePayload struct {
	MessageID api.Snowflake  `json:"id"`
	ChannelID api.Snowflake  `json:"channel_id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
}

// MessageDeleteHandler handles api.GatewayEventMessageDelete
type MessageDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h *MessageDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageDeleteHandler) New() interface{} {
	return &messageDeletePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageDeletePayload)
	if !ok {
		return
	}

	message := disgo.Cache().Message(payload.ChannelID, payload.MessageID)
	disgo.Cache().UncacheMessage(payload.ChannelID, payload.MessageID)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		ChannelID:    payload.ChannelID,
		Message:      message,
	}

	eventManager.Dispatch(&events.MessageDeleteEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if payload.GuildID == nil {
		eventManager.Dispatch(&events.DMMessageDeleteEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		eventManager.Dispatch(&events.GuildMessageDeleteEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *payload.GuildID,
			},
		})
	}
}
