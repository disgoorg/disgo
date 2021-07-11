package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// CommandUpdateHandler handles api.CommandCreateEvent
type CommandUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h CommandUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventCommandUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h CommandUpdateHandler) New() interface{} {
	return &api.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h CommandUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*api.Command)
	if !ok {
		return
	}

	var oldCommand *api.Command
	if command.ApplicationID == disgo.ApplicationID() {
		oldCommand = disgo.Cache().Command(command.ID)
		if oldCommand != nil {
			oldCommand = &*oldCommand
		}
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, api.CacheStrategyYes)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, api.CacheStrategyYes)
	}

	eventManager.Dispatch(&events.CommandUpdateEvent{
		GenericCommandEvent: &events.GenericCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
		// always nil for not our own commands
		OldCommand: oldCommand,
	})
}
