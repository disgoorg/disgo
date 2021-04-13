package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericApplicationCommandEvent is called upon receiving either ApplicationCommandCreateEvent, ApplicationCommandUpdateEvent or ApplicationCommandDeleteEvent
type GenericApplicationCommandEvent struct {
	GenericEvent
	CommandID api.Snowflake
	Command   *api.Command
	GuildID   *api.Snowflake
}

// Guild returns the api.Guild the api.Event got called or nil for global api.Command(s)
func (e GenericApplicationCommandEvent) Guild() *api.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().Guild(*e.GuildID)
}

// ApplicationCommandCreateEvent indicates that a new api.Command got created(this can come from any bot!)
type ApplicationCommandCreateEvent struct {
	GenericApplicationCommandEvent
}

// ApplicationCommandUpdateEvent indicates that a api.Command got updated(this can come from any bot!)
type ApplicationCommandUpdateEvent struct {
	GenericApplicationCommandEvent
	OldCommand *api.Command
}

// ApplicationCommandDeleteEvent indicates that a api.Command got deleted(this can come from any bot!)
type ApplicationCommandDeleteEvent struct {
	GenericApplicationCommandEvent
}
