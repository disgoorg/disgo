package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMChannelUserTypingEvent
type GenericDMChannelEvent struct {
	*GenericChannelEvent
}

// DMChannelCreateEvent indicates that a new api.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that an api.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldChannel *Channel
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

// DMChannelUserTypingEvent indicates that an api.User started typing in an api.DMChannel(requires api.GatewayIntentsDirectMessageTyping)
type DMChannelUserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// DMChannel returns the core.DMChannel the DMChannelUserTypingEvent happened in
func (e DMChannelUserTypingEvent) DMChannel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
