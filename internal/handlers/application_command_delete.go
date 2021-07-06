package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// CommandDeleteHandler handles api.CommandCreateEvent
type CommandDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h CommandDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventCommandDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h CommandDeleteHandler) New() interface{} {
	return &api.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h CommandDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*api.Command)
	if !ok {
		return
	}

	// we only cache our own commands
	if command.ApplicationID == disgo.ApplicationID() {
		disgo.Cache().UncacheCommand(command.ID)
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, api.CacheStrategyNo)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, api.CacheStrategyNo)
	}

	eventManager.Dispatch(&events.CommandDeleteEvent{
		GenericCommandEvent: &events.GenericCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
	})
}
