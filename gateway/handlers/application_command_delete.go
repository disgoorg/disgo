package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// CommandDeleteHandler handles api.ApplicationCommandCreateEvent
type CommandDeleteHandler struct{}

// EventType returns the api.GatewayEventType
func (h *CommandDeleteHandler) EventType() gateway.EventType {
	return gateway.EventTypeCommandDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *CommandDeleteHandler) New() interface{} {
	return discord.ApplicationCommand{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *CommandDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(discord.ApplicationCommand)
	if !ok {
		return
	}

	// we only cache our own commands
	if command.ApplicationID == disgo.ApplicationID() {
		disgo.Cache().UncacheCommand(command.ID)
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, core.CacheStrategyNo)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, core.CacheStrategyNo)
	}

	eventManager.Dispatch(&events.ApplicationCommandDeleteEvent{
		GenericApplicationCommandEvent: &events.GenericApplicationCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
	})
}
