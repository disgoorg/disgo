package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMMessageEvent generic api.DMChannel api.Message api.Event
type GenericDMMessageEvent struct {
	GenericDMEvent
	GenericMessageEvent
}

// DMMessageReceivedEvent called upon receiving a api.Message in a api.DMChannel
type DMMessageReceivedEvent struct {
	GenericDMMessageEvent
	Message *api.Message
}

// DMMessageUpdateEvent called upon editing a api.Message in a api.DMChannel
type DMMessageUpdateEvent struct {
	GenericDMMessageEvent
	Message *api.Message
}

// DMMessageDeleteEvent called upon deleting a api.Message in a api.DMChannel
type DMMessageDeleteEvent struct {
	GenericDMMessageEvent
}
