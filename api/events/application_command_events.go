package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericApplicationCommandEvent struct {
	GenericEvent
	CommandID api.Snowflake
	GuildID   *api.Snowflake
}

func (e GenericApplicationCommandEvent) Guild() *api.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().Guild(*e.GuildID)
}

type ApplicationCommandCreateEvent struct {
	GenericApplicationCommandEvent
	Command *api.Command
}

type ApplicationCommandUpdateEvent struct {
	GenericApplicationCommandEvent
	NewCommand *api.Command
	OldCommand *api.Command
}

type ApplicationCommandDeleteEvent struct {
	GenericApplicationCommandEvent
	Command *api.Command
}
