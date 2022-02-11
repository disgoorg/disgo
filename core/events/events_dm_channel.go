package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMUserTypingStartEvent
type GenericDMChannelEvent struct {
	*GenericEvent
	Channel   *core.DMChannel
	ChannelID snowflake.Snowflake
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
	ChannelID           snowflake.Snowflake
	NewLastPinTimestamp *discord.Time
	OldLastPinTimestamp *discord.Time
}

// DMUserTypingStartEvent indicates that a core.User started typing in a core.DMChannel(requires discord.GatewayIntentDirectMessageTyping)
type DMUserTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	UserID    snowflake.Snowflake
	Timestamp time.Time
}

// Channel returns the core.DMChannel the DMUserTypingStartEvent happened in
func (e DMUserTypingStartEvent) Channel() *core.DMChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}
