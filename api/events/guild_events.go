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

type GuildUpdateEvent struct {
	GenericGuildEvent
	Guild    *api.Guild
	OldGuild *api.Guild
}

type GuildAvailableEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

type GuildUnavailableEvent struct {
	GenericGuildEvent
}

type GuildJoinEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

type GuildLeaveEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}
