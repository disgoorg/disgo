package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMUserTypingEvent
type GenericDMChannelEvent struct {
	*GenericChannelEvent
	DMChannel core.DMChannel
}

// DMChannelCreateEvent indicates that a new api.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that an api.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldDMChannel core.DMChannel
}

type DMChannelPinsUpdateEvent struct {
	*GenericDMChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// DMChannelDeleteEvent indicates that an api.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

// DMUserTypingEvent indicates that an api.User started typing in an api.DMChannel(requires api.GatewayIntentsDirectMessageTyping)
type DMUserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// DMChannel returns the core.DMChannel the DMUserTypingEvent happened in
func (e DMUserTypingEvent) DMChannel() core.DMChannel {
	return e.Disgo().Cache().DMChannelCache().Get(e.ChannelID)
}
