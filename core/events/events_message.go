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
	GuildID   *discord.Snowflake
}

// Channel returns the core.Channel where the GenericMessageEvent happened
func (e *GenericMessageEvent) Channel() core.MessageChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.MessageChannel)
	}
	return nil
}

// Guild returns the core.Guild where the GenericMessageEvent happened or nil if it happened in DMs
func (e *GenericMessageEvent) Guild() *core.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Bot().Caches.Guilds().Get(*e.GuildID)
}

// MessageCreateEvent indicates that a core.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that a core.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *core.Message
}

// MessageDeleteEvent indicates that a core.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
