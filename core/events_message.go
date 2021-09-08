package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericMessageEvent generic api.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *Message
	ChannelID discord.Snowflake
}

// Channel returns the core.Channel where the GenericMessageEvent happened
func (e *GenericMessageEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// MessageCreateEvent indicates that an api.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that an api.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *Message
}

// MessageDeleteEvent indicates that an api.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
