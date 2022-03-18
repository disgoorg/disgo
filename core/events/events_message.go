package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericMessageEvent generic core.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID snowflake.Snowflake
	Message   discord.Message
	ChannelID snowflake.Snowflake
	GuildID   *snowflake.Snowflake
}

// Channel returns the core.Channel where the GenericMessageEvent happened
func (e *GenericMessageEvent) Channel() (discord.MessageChannel, bool) {
	return e.Bot().Caches().Channels().GetMessageChannel(e.ChannelID)
}

// Guild returns the core.Guild where the GenericMessageEvent happened or nil if it happened in DMs
func (e *GenericMessageEvent) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Bot().Caches().Guilds().Get(*e.GuildID)
}

// MessageCreateEvent indicates that a core.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that a core.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage discord.Message
}

// MessageDeleteEvent indicates that a core.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
