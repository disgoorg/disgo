package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericMessageEvent generic core.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *core.Message
	ChannelID discord.Snowflake
}

// Channel returns the core.Channel where the GenericMessageEvent happened
func (e *GenericMessageEvent) Channel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// MessageCreateEvent indicates that an core.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that an core.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *core.Message
}

// MessageDeleteEvent indicates that an core.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
