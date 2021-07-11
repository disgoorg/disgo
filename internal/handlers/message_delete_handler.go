package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

type messageDeletePayload struct {
	MessageID api.Snowflake  `json:"id"`
	GuildID   *api.Snowflake `json:"guild_id,omitempty"`
	ChannelID api.Snowflake  `json:"channel_id"`
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
func (h *MessageDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	//payload, ok := i.(*api.messageDeletePayload)
}
