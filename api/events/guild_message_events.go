package events

import (
	"github.com/DiscoOrg/disgo/api"
)

// GenericGuildMessageEvent indicates that we received a api.Message api.Event in a api.Guild
type GenericGuildMessageEvent struct {
	GenericGuildEvent
	GenericMessageEvent
}

// GuildMessageReceivedEvent indicates that we received a api.Message in a api.Guild
type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}

// GuildMessageUpdateEvent indicates that a api.Message was updated in a api.Guild
type GuildMessageUpdateEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}

// GuildMessageDeleteEvent indicates that a api.Message was deleted in a api.Guild
type GuildMessageDeleteEvent struct {
	GenericGuildMessageEvent
	Message *api.Message
}
