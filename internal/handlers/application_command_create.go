package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// ApplicationCommandCreateHandler handles api.ApplicationCommandCreateEvent
type ApplicationCommandCreateHandler struct{}

// Event returns the raw gateway event Event
func (h ApplicationCommandCreateHandler) Event() api.GatewayEventName {
	return api.GatewayEventApplicationCommandCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h ApplicationCommandCreateHandler) New() interface{} {
	return &api.Command{}
}

// Handle handles the specific raw gateway event
func (h ApplicationCommandCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	/*command, ok := i.(*api.Command)
	if !ok {
		return
	}*/
}
