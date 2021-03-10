package events

import (
	"github.com/DiscoOrg/disgo/api"
)

type GenericMessageEvent struct {
	api.Event
	MessageID api.Snowflake
}

type GenericGuildMessageEvent struct {
	GenericMessageEvent
	GenericGuildEvent
	TextChannel api.TextChannel
}


type MessageReceivedEvent struct {
	GenericMessageEvent
	Message api.Message
}

type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}
