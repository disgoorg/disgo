package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMChannelUserTypingStartEvent
type GenericDMChannelEvent struct {
	*GenericEvent
	Channel   *core.DMChannel
	ChannelID discord.Snowflake
}

// DMChannelCreateEvent indicates that a new core.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that a core.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldChannel *core.DMChannel
}

// DMChannelDeleteEvent indicates that a core.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

type DMChannelPinsUpdateEvent struct {
	*GenericEvent
	ChannelID           discord.Snowflake
	NewLastPinTimestamp *discord.Time
	OldLastPinTimestamp *discord.Time
}

// DMChannelUserTypingStartEvent indicates that a core.User started typing in a core.DMChannel(requires discord.GatewayIntentDirectMessageTyping)
type DMChannelUserTypingStartEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	UserID    discord.Snowflake
	Timestamp time.Time
}

// Channel returns the core.DMChannel the DMChannelUserTypingStartEvent happened in
func (e DMChannelUserTypingStartEvent) Channel() *core.DMChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}
