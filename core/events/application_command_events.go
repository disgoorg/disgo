package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// GenericCommandEvent is called upon receiving either CommandCreateEvent, CommandUpdateEvent or CommandDeleteEvent
type GenericCommandEvent struct {
	*GenericEvent
	Command *core.Command
}

// Guild returns the api.Guild the api.EventType got called or nil for global api.Command(s)
func (e GenericCommandEvent) Guild() *core.Guild {
	if e.Command.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().GuildCache().Get(*e.Command.GuildID)
}

// CommandCreateEvent indicates that a new api.Command got created(this can come from any bot!)
type CommandCreateEvent struct {
	*GenericCommandEvent
}

// CommandUpdateEvent indicates that an api.Command got updated(this can come from any bot!)
type CommandUpdateEvent struct {
	*GenericCommandEvent
	OldCommand *core.Command
}

// CommandDeleteEvent indicates that an api.Command got deleted(this can come from any bot!)
type CommandDeleteEvent struct {
	*GenericCommandEvent
}
