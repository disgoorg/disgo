package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericCommandEvent is called upon receiving either CommandCreateEvent, CommandUpdateEvent or CommandDeleteEvent
type GenericCommandEvent struct {
	GenericEvent
	CommandID api.Snowflake
	Command   *api.Command
	GuildID   *api.Snowflake
}

// Guild returns the api.Guild the api.Event got called or nil for global api.Command(s)
func (e GenericCommandEvent) Guild() *api.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().Guild(*e.GuildID)
}

// CommandCreateEvent indicates that a new api.Command got created(this can come from any bot!)
type CommandCreateEvent struct {
	GenericCommandEvent
}

// CommandUpdateEvent indicates that a api.Command got updated(this can come from any bot!)
type CommandUpdateEvent struct {
	GenericCommandEvent
	OldCommand *api.Command
}

// CommandDeleteEvent indicates that a api.Command got deleted(this can come from any bot!)
type CommandDeleteEvent struct {
	GenericCommandEvent
}
