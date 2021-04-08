package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericDMChannelEvent struct {
	GenericChannelEvent
}

func (e GenericDMChannelEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}

type DMChannelCreateEvent struct {
	GenericDMChannelEvent
	DMChannel *api.DMChannel
}

type DMChannelUpdateEvent struct {
	GenericDMChannelEvent
	NewDMChannel *api.DMChannel
	OldDMChannel *api.DMChannel
}

type DMChannelDeleteEvent struct {
	GenericDMChannelEvent
	DMChannel *api.DMChannel
}
