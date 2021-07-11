package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMUserTypingEvent
type GenericDMChannelEvent struct {
	*GenericChannelEvent
	DMChannel *api.DMChannel
}

// DMChannelCreateEvent indicates that a new api.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that a api.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldDMChannel *api.DMChannel
}

// DMChannelDeleteEvent indicates that a api.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

// DMUserTypingEvent indicates that a api.User started typing in a api.DMChannel(requires api.GatewayIntentsDirectMessageTyping)
type DMUserTypingEvent struct {
	*GenericDMChannelEvent
}
