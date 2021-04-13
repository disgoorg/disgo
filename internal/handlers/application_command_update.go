package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
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

// HandleGatewayEvent handles the specific raw gateway event
func (h ApplicationCommandUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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

	genericApplicationCommandEvent := events.GenericApplicationCommandEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		CommandID:    command.ID,
		Command:      command,
		GuildID:      command.GuildID,
	}
	eventManager.Dispatch(genericApplicationCommandEvent)

	eventManager.Dispatch(events.ApplicationCommandUpdateEvent{
		GenericApplicationCommandEvent: genericApplicationCommandEvent,
		// always nil for not our own commands
		OldCommand: oldCommand,
	})
}
