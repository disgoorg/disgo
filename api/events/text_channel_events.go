package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericTextChannelEvent struct {
	GenericChannelEvent
	TextChannel *api.TextChannel
}

type TextChannelCreateEvent struct {
	GenericTextChannelEvent
}

type TextChannelUpdateEvent struct {
	GenericTextChannelEvent
	OldTextChannel *api.TextChannel
}

type TextChannelDeleteEvent struct {
	GenericTextChannelEvent
}
