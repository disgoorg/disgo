package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// DMChannelPinsUpdate indicates that a discord.Message got pinned or unpinned.
type DMChannelPinsUpdate struct {
	*GenericEvent
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
}

// DMUserTypingStart indicates that a discord.User started typing in a discord.DMChannel(requires gateway.IntentDirectMessageTyping).
type DMUserTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}
