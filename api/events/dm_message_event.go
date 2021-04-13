package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmoteEvent or DMMessageReactionRemoveAllEvent(requires api.IntentsDirectMessages)
type GenericDMMessageEvent struct {
	GenericMessageEvent
	Message *api.Message
}

// DMChannel returns the api.DMChannel where the GenericDMMessageEvent happened
func (e GenericDMMessageEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}

// DMMessageCreateEvent is called upon receiving a api.Message in a api.DMChannel(requires api.IntentsDirectMessages)
type DMMessageCreateEvent struct {
	GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing a api.Message in a api.DMChannel(requires api.IntentsDirectMessages)
type DMMessageUpdateEvent struct {
	GenericDMMessageEvent
	OldMessage *api.Message
}

// DMMessageDeleteEvent is called upon deleting a api.Message in a api.DMChannel(requires api.IntentsDirectMessages)
type DMMessageDeleteEvent struct {
	GenericDMMessageEvent
}
