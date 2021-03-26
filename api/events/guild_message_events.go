package events

import (
	"github.com/DiscoOrg/disgo/api"
)

type GenericGuildMessageEvent struct {
	GenericGuildEvent
	GenericMessageEvent
}

type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}

type GuildMessageUpdateEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}
