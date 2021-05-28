package handlers

import "github.com/DisgoOrg/disgo/api"

// MessageUpdateHandler handles api.GatewayEventMessageUpdate
type MessageUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h MessageUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageUpdateHandler) New() interface{} {
	return &api.FullMessage{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	//fullMessage, ok := i.(*api.FullMessage)
}
