package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMMessageEvent generic api.DMChannel api.Message api.GenericEvent
type GenericDMMessageEvent struct {
	GenericMessageEvent
	Message *api.Message
}

func (e GenericDMMessageEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}

// DMMessageReceivedEvent called upon receiving a api.Message in a api.DMChannel
type DMMessageReceivedEvent struct {
	GenericDMMessageEvent
}

// DMMessageUpdateEvent called upon editing a api.Message in a api.DMChannel
type DMMessageUpdateEvent struct {
	GenericDMMessageEvent
	OldMessage *api.Message
}

// DMMessageDeleteEvent called upon deleting a api.Message in a api.DMChannel
type DMMessageDeleteEvent struct {
	GenericDMMessageEvent
}
