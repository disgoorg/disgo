package events

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMChannelEvent is called upon receiving DMChannelCreateEvent, DMChannelUpdateEvent, DMChannelDeleteEvent or DMUserTypingStartEvent
type GenericDMChannelEvent struct {
	*GenericEvent
	Channel   discord.DMChannel
	ChannelID snowflake.ID
}

// DMChannelCreateEvent indicates that a new discord.DMChannel got created
type DMChannelCreateEvent struct {
	*GenericDMChannelEvent
}

// DMChannelUpdateEvent indicates that a discord.DMChannel got updated
type DMChannelUpdateEvent struct {
	*GenericDMChannelEvent
	OldChannel discord.DMChannel
}

// DMChannelDeleteEvent indicates that a discord.DMChannel got deleted
type DMChannelDeleteEvent struct {
	*GenericDMChannelEvent
}

type DMChannelPinsUpdateEvent struct {
	*GenericEvent
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
	OldLastPinTimestamp *time.Time
}

// DMUserTypingStartEvent indicates that a discord.User started typing in a discord.DMChannel(requires discord.GatewayIntentDirectMessageTyping)
type DMUserTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}

// Channel returns the discord.DMChannel the DMUserTypingStartEvent happened in
func (e DMUserTypingStartEvent) Channel() (discord.DMChannel, bool) {
	if channel, ok := e.Client().Caches().Channels().Get(e.ChannelID); ok {
		return channel.(discord.DMChannel), false
	}
	return discord.DMChannel{}, true
}
