package events

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/models"
)

type GenericGuildEvent struct {
	api.Event
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
