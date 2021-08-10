package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

type messageDeletePayload struct {
	MessageID discord.Snowflake  `json:"id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	ChannelID discord.Snowflake  `json:"channel_id"`
}

// MessageDeleteHandler handles api.GatewayEventMessageDelete
type MessageDeleteHandler struct{}

// Event returns the api.GatewayEventType
func (h *MessageDeleteHandler) EventType() gateway.EventType {
	return gateway.EventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageDeleteHandler) New() interface{} {
	return &messageDeletePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageDeletePayload)
	if !ok {
		return
	}

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		Message:      disgo.Cache().Message(payload.ChannelID, payload.MessageID),
		ChannelID:    payload.ChannelID,
	}

	disgo.EventManager().Dispatch(&events.MessageDeleteEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if payload.GuildID == nil {
		disgo.EventManager().Dispatch(&events.DMMessageDeleteEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		disgo.EventManager().Dispatch(&events.GuildMessageDeleteEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *payload.GuildID,
			},
		})
	}
}
