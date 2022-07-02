package events

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMChannel is called upon receiving DMChannelCreate , DMChannelUpdate , DMChannelDelete or DMUserTypingStart.
type GenericDMChannel struct {
	*GenericEvent
	Channel   discord.DMChannel
	ChannelID snowflake.ID
}

// DMChannelCreate indicates that a new discord.DMChannel got created.
type DMChannelCreate struct {
	*GenericDMChannel
}

// DMChannelUpdate indicates that a discord.DMChannel got updated.
type DMChannelUpdate struct {
	*GenericDMChannel
	OldChannel discord.DMChannel
}

// DMChannelDelete indicates that a discord.DMChannel got deleted.
type DMChannelDelete struct {
	*GenericDMChannel
}

// DMChannelPinsUpdate indicates that a discord.Message got pinned or unpinned.
type DMChannelPinsUpdate struct {
	*GenericEvent
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
	OldLastPinTimestamp *time.Time
}

// DMUserTypingStart indicates that a discord.User started typing in a discord.DMChannel(requires gateway.IntentDirectMessageTyping).
type DMUserTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}

// Channel returns the discord.DMChannel the DMUserTypingStart happened in.
func (e DMUserTypingStart) Channel() (discord.DMChannel, bool) {
	if channel, ok := e.Client().Caches().Channels().Get(e.ChannelID); ok {
		return channel.(discord.DMChannel), false
	}
	return discord.DMChannel{}, true
}
