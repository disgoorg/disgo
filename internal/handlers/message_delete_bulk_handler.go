package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type messageDeleteBulkPayload struct {
	MessageIDs []api.Snowflake `json:"ids"`
	ChannelID  api.Snowflake   `json:"channel_id"`
	GuildID    *api.Snowflake  `json:"guild_id,omitempty"`
}

// MessageDeleteBulkHandler handles api.GatewayEventMessageBulkDelete
type MessageDeleteBulkHandler struct{}

// Event returns the raw gateway event Event
func (h MessageDeleteBulkHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageDeleteBulk
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageDeleteBulkHandler) New() interface{} {
	return &messageDeleteBulkPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageDeleteBulkHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*messageDeleteBulkPayload)
	if !ok {
		return
	}

	for _, messageID := range payload.MessageIDs {
		message := disgo.Cache().Message(payload.ChannelID, messageID)
		disgo.Cache().UncacheMessage(payload.ChannelID, messageID)

		genericMessageEvent := &events.GenericMessageEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			MessageID:    messageID,
			ChannelID:    payload.ChannelID,
			Message:      message,
		}

		eventManager.Dispatch(&events.MessageDeleteEvent{
			GenericMessageEvent: genericMessageEvent,
		})

		if message.GuildID == nil {
			eventManager.Dispatch(&events.DMMessageDeleteEvent{
				GenericDMMessageEvent: &events.GenericDMMessageEvent{
					GenericMessageEvent: genericMessageEvent,
				},
			})
		} else {
			eventManager.Dispatch(&events.GuildMessageDeleteEvent{
				GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
					GenericMessageEvent: genericMessageEvent,
					GuildID:             *message.GuildID,
				},
			})
		}
	}

}
