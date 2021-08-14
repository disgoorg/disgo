package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// GenericApplicationCommandEvent is called upon receiving either ApplicationCommandCreateEvent, ApplicationCommandUpdateEvent or ApplicationCommandDeleteEvent
type GenericApplicationCommandEvent struct {
	*GenericEvent
	Command *core.ApplicationCommand
}

// Guild returns the api.Guild the api.EventType got called or nil for global api.Command(s)
func (e *GenericApplicationCommandEvent) Guild() *core.Guild {
	if e.Command.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().GuildCache().Get(*e.Command.GuildID)
}

// ApplicationCommandCreateEvent indicates that a new api.Command got created(this can come from any bot!)
type ApplicationCommandCreateEvent struct {
	*GenericApplicationCommandEvent
}

// ApplicationCommandUpdateEvent indicates that an api.Command got updated(this can come from any bot!)
type ApplicationCommandUpdateEvent struct {
	*GenericApplicationCommandEvent
	OldCommand *core.ApplicationCommand
}

// ApplicationCommandDeleteEvent indicates that an api.Command got deleted(this can come from any bot!)
type ApplicationCommandDeleteEvent struct {
	*GenericApplicationCommandEvent
}
