package events

import (
	"github.com/DiscoOrg/disgo/api"
)

type GenericGuildEvent struct {
	api.Event
	Guild api.Guild
}

type GuildAvailableEvent struct {
	GenericGuildEvent
}

type GuildUnavailableEvent struct {
	GenericGuildEvent
}

type GuildJoinEvent struct {
	GenericGuildEvent
}

type GuildLeaveEvent struct {
	GenericGuildEvent
}
