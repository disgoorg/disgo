package events

import (
	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/models"
)

type GenericGuildEvent struct {
	disgo.Event
	Guild models.Guild
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
