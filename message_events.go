package disgo

import (
	"github.com/DiscoOrg/disgo/models"
)

type GenericMessageEvent struct {
	Event
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
