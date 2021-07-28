package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// CommandCreateHandler handles api.CommandCreateEvent
type CommandCreateHandler struct{}

// Event returns the raw gateway event Event
func (h *CommandCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventCommandCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *CommandCreateHandler) New() interface{} {
	return &api.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *CommandCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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

	eventManager.Dispatch(&events.CommandCreateEvent{
		GenericCommandEvent: &events.GenericCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
	})
}
