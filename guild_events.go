package disgo

import (
	"github.com/DiscoOrg/disgo/models"
)

type GenericGuildEvent struct {
	Event
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
