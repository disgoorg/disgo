package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// ChannelUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h ChannelUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h ChannelUpdateHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.Channel)
	if !ok {
		return
	}

	
}
