package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// ApplicationCommandDeleteHandler handles api.ApplicationCommandCreateEvent
type ApplicationCommandDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h ApplicationCommandDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventApplicationCommandDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h ApplicationCommandDeleteHandler) New() interface{} {
	return &api.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ApplicationCommandDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*api.Command)
	if !ok {
		return
	}

	genericApplicationCommandEvent := events.GenericApplicationCommandEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		CommandID:    command.ID,
		GuildID:      command.GuildID,
	}
	eventManager.Dispatch(genericApplicationCommandEvent)

	eventManager.Dispatch(events.ApplicationCommandDeleteEvent{
		GenericApplicationCommandEvent: genericApplicationCommandEvent,
		Command: command,
	})
}
