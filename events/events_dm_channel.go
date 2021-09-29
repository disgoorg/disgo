package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMChannelUserTypingEvent
type GenericDMChannelEvent struct {
	*GenericChannelEvent
}

// DMChannelCreateEvent indicates that a Channel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that a Channel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldChannel *core.Channel
}

type DMChannelPinsUpdateEvent struct {
	*GenericDMChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// DMChannelDeleteEvent indicates that a Channel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUserTypingEvent indicates that a core.User started typing in a Channel (requires discord.GatewayIntentDirectMessageTyping)
type DMChannelUserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// DMChannel returns the Channel the DMChannelUserTypingEvent happened in.
// This will only check cached channels!
func (e DMChannelUserTypingEvent) DMChannel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
