package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

// CommandCreateHandler handles api.ApplicationCommandCreateEvent
type CommandCreateHandler struct{}

// EventType returns the raw gateway api.GatewayEventType
func (h *CommandCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeCommandCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *CommandCreateHandler) New() interface{} {
	return &discord.ApplicationCommand{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *CommandCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	command, ok := i.(*discord.ApplicationCommand)
	if !ok {
		return
	}

	// only cache our own commands
	cacheStrategy := core.CacheStrategyNo
	if command.ApplicationID == disgo.ApplicationID() {
		cacheStrategy = core.CacheStrategyYes
	}

	if command.FromGuild() {
		command = disgo.EntityBuilder().CreateGuildCommand(*command.GuildID, command, cacheStrategy)
	} else {
		command = disgo.EntityBuilder().CreateGlobalCommand(command, cacheStrategy)
	}

	eventManager.Dispatch(&events.ApplicationCommandCreateEvent{
		GenericApplicationCommandEvent: &events.GenericApplicationCommandEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			Command:      command,
		},
	})
}
