package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type messageDeletePayload struct {
	MessageID discord.Snowflake  `json:"id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	ChannelID discord.Snowflake  `json:"channel_id"`
}

// MessageDeleteHandler handles api.GatewayEventMessageDelete
type MessageDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *MessageDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageDeleteHandler) New() interface{} {
	return messageDeletePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(messageDeletePayload)
	if !ok {
		return
	}

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    payload.MessageID,
		Message:      disgo.Cache().MessageCache().GetCopy(payload.ChannelID, payload.MessageID),
		ChannelID:    payload.ChannelID,
	}

	disgo.Cache().MessageCache().Uncache(payload.ChannelID, payload.MessageID)

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
