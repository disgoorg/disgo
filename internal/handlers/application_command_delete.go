package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// ApplicationCommandDeleteHandler handles api.ApplicationCommandCreateEvent
type ApplicationCommandDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h ApplicationCommandDeleteHandler) Event() api.GatewayEventName {
	return api.GatewayEventApplicationCommandDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h ApplicationCommandDeleteHandler) New() interface{} {
	return &api.Command{}
}

// Handle handles the specific raw gateway event
func (h ApplicationCommandDeleteHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	/*command, ok := i.(*api.Command)
	if !ok {
		return
	}*/
}
