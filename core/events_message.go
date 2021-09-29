package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericMessageEvent generic core.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *Message
	ChannelID discord.Snowflake
}

// Channel returns the core.Channel the GenericMessageEvent happened in.
// This will only check cached channels!
func (e *GenericMessageEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// MessageCreateEvent indicates that a core.Message got created
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that a core.Message got updated
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *Message
}

// MessageDeleteEvent indicates that a core.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
