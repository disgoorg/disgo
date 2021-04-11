package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// ApplicationCommandCreateHandler handles api.ApplicationCommandCreateEvent
type ApplicationCommandCreateHandler struct{}

// Event returns the raw gateway event Event
func (h ApplicationCommandCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventApplicationCommandCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h ApplicationCommandCreateHandler) New() interface{} {
	return &api.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ApplicationCommandCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*api.Command)
	if !ok {
		return
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, api.CacheStrategyYes)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, api.CacheStrategyYes)
	}

	genericApplicationCommandEvent := events.GenericApplicationCommandEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		CommandID:    command.ID,
		GuildID:      command.GuildID,
	}
	eventManager.Dispatch(genericApplicationCommandEvent)

	eventManager.Dispatch(events.ApplicationCommandCreateEvent{
		GenericApplicationCommandEvent: genericApplicationCommandEvent,
		Command: command,
	})
}
