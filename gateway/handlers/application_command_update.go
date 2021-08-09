package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

// CommandUpdateHandler handles api.CommandCreateEvent
type CommandUpdateHandler struct{}

// Event returns the api.GatewayEventType
func (h *CommandUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeCommandUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *CommandUpdateHandler) New() interface{} {
	return &discord.Command{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *CommandUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*discord.Command)
	if !ok {
		return
	}

	var oldCommand *discord.Command
	if command.ApplicationID == disgo.ApplicationID() {
		oldCommand = disgo.Cache().Command(command.ID)
		if oldCommand != nil {
			oldCommand = &*oldCommand
		}
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, core.CacheStrategyYes)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, core.CacheStrategyYes)
	}

	eventManager.Dispatch(&events.CommandUpdateEvent{
		GenericCommandEvent: &events.GenericCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
		// this is always nil for not our own commands
		OldCommand: oldCommand,
	})
}
