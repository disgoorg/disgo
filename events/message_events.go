package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericMessageEvent generic discord.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID snowflake.Snowflake
	Message   discord.Message
	ChannelID snowflake.Snowflake
	GuildID   *snowflake.Snowflake
}

// Channel returns the discord.Channel where the GenericMessageEvent happened
func (e *GenericMessageEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID)
}

// Guild returns the discord.Guild where the GenericMessageEvent happened or nil if it happened in DMs
func (e *GenericMessageEvent) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Client().Caches().Guilds().Get(*e.GuildID)
}

// MessageCreateEvent indicates that a discord.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that a discord.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage discord.Message
}

// MessageDeleteEvent indicates that a discord.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}
