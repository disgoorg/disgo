package events

import (
	"github.com/DiscoOrg/disgo/api"
)

type GenericDMMessageEvent struct {
	GenericDMEvent
	GenericMessageEvent
}


type DMMessageReceivedEvent struct {
	GenericDMMessageEvent
	Message api.Message
}


type DMMessageUpdateEvent struct {
	GenericDMMessageEvent
	Message api.Message
}

type DMMessageDeleteEvent struct {
	GenericDMMessageEvent
}
