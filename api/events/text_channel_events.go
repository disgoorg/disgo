package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericTextChannelEvent struct {
	GenericChannelEvent
}

func (e GenericTextChannelEvent) Category() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

type TextChannelCreateEvent struct {
	GenericTextChannelEvent
	TextChannel *api.TextChannel
}

type TextChannelUpdateEvent struct {
	GenericTextChannelEvent
	NewTextChannel *api.TextChannel
	OldTextChannel *api.TextChannel
}

type TextChannelDeleteEvent struct {
	GenericTextChannelEvent
	TextChannel *api.TextChannel
}
