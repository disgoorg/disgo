package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMChannelUserTypingEvent
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
	OldChannel core.Channel
}

// DMChannelDeleteEvent indicates that a core.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUserTypingEvent indicates that a core.User started typing in a core.DMChannel(requires discord.GatewayIntentDirectMessageTyping)
type DMChannelUserTypingEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	UserID    discord.Snowflake
	Timestamp time.Time
}

// DMChannel returns the core.DMChannel the DMChannelUserTypingEvent happened in
func (e DMChannelUserTypingEvent) DMChannel() *core.DMChannel {
	if ch := e.Bot().Caches.ChannelCache().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}
