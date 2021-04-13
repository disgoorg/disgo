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

	// only cache our own commands
	cacheStrategy := api.CacheStrategyNo
	if command.ApplicationID == disgo.ApplicationID() {
		cacheStrategy = api.CacheStrategyYes
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, cacheStrategy)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, cacheStrategy)
	}

	genericApplicationCommandEvent := events.GenericApplicationCommandEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		CommandID:    command.ID,
		Command:      command,
		GuildID:      command.GuildID,
	}
	eventManager.Dispatch(genericApplicationCommandEvent)

	eventManager.Dispatch(events.ApplicationCommandCreateEvent{
		GenericApplicationCommandEvent: genericApplicationCommandEvent,
	})
}
