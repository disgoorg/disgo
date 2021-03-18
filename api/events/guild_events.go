package events

import (
	"github.com/DiscoOrg/disgo/api"
)

type GenericGuildEvent struct {
	api.Event
	GuildID api.Snowflake
}
func (e GenericGuildEvent) Guild() *api.Guild {
	return e.Disgo.Cache().Guild(e.GuildID)
}


type GuildAvailableEvent struct {
	GenericGuildEvent
}

type GuildUnavailableEvent struct {
	GenericGuildEvent
	Unavailable bool
}

type GuildJoinEvent struct {
	GenericGuildEvent
}

type GuildLeaveEvent struct {
	GenericGuildEvent
}
