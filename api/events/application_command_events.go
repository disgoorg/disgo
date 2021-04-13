package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericApplicationCommandEvent struct {
	GenericEvent
	CommandID api.Snowflake
	Command   *api.Command
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
}

type ApplicationCommandUpdateEvent struct {
	GenericApplicationCommandEvent
	OldCommand *api.Command
}

type ApplicationCommandDeleteEvent struct {
	GenericApplicationCommandEvent
}
