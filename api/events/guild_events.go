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
	Guild *api.Guild
}

type GuildUnavailableEvent struct {
	GenericGuildEvent
	Unavailable bool
}

type GuildJoinEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

type GuildLeaveEvent struct {
	GenericGuildEvent
}
