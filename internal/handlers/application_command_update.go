package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// ApplicationCommandUpdateHandler handles api.ApplicationCommandCreateEvent
type ApplicationCommandUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h ApplicationCommandUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventApplicationCommandUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h ApplicationCommandUpdateHandler) New() interface{} {
	return &api.Command{}
}

// Handle handles the specific raw gateway event
func (h ApplicationCommandUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	/*command, ok := i.(*api.Command)
	if !ok {
		return
	}*/
}
