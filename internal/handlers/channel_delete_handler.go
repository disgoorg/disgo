package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h ChannelDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h ChannelDeleteHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.Channel)
	if !ok {
		return
	}

	
}
