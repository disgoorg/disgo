package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericApplicationCommandEvent struct {
	GenericEvent
	CommandID api.Snowflake
	GuildID   *api.Snowflake
	Guild     *api.Guild
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
