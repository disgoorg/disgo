package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMChannelUserTypingEvent
type GenericDMChannelEvent struct {
	*GenericChannelEvent
}

// DMChannelCreateEvent indicates that a new core.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that an core.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldChannel core.Channel
}

type DMChannelPinsUpdateEvent struct {
	*GenericDMChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// DMChannelDeleteEvent indicates that an core.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUserTypingEvent indicates that an core.User started typing in an core.DMChannel(requires core.GatewayIntentsDirectMessageTyping)
type DMChannelUserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// DMChannel returns the core.DMChannel the DMChannelUserTypingEvent happened in
func (e DMChannelUserTypingEvent) DMChannel() *core.DMChannel {
	if ch := e.Bot().Caches.ChannelCache().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}
