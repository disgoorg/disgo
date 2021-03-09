package events

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/models"
)

type GenericMessageEvent struct {
	api.Event
	MessageID models.Snowflake
}

type GenericGuildMessageEvent struct {
	GenericMessageEvent
	GenericGuildEvent
	TextChannel models.TextChannel
}


type MessageReceivedEvent struct {
	GenericMessageEvent
	Message models.Message
}

type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message models.Message
}
