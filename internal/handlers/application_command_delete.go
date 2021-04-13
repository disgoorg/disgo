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

	// we only cache our own commands
	if command.ApplicationID == disgo.ApplicationID() {
		disgo.Cache().UncacheCommand(command.ID)
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, api.CacheStrategyNo)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, api.CacheStrategyNo)
	}

	genericApplicationCommandEvent := events.GenericApplicationCommandEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		CommandID:    command.ID,
		Command:      command,
		GuildID:      command.GuildID,
	}
	eventManager.Dispatch(genericApplicationCommandEvent)

	eventManager.Dispatch(events.ApplicationCommandDeleteEvent{
		GenericApplicationCommandEvent: genericApplicationCommandEvent,
	})
}
