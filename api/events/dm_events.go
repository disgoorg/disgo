package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericDMChannelEvent struct {
	GenericChannelEvent
	DMChannel *api.DMChannel
}

type DMChannelCreateEvent struct {
	GenericDMChannelEvent
}

type DMChannelUpdateEvent struct {
	GenericDMChannelEvent
	OldDMChannel *api.DMChannel
}

type DMChannelDeleteEvent struct {
	GenericDMChannelEvent
}

type DMUserTypingEvent struct {
	GenericDMChannelEvent
}

func (e DMUserTypingEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}
