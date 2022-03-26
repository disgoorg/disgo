package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake"
)

// GenericGuildMessageEvent is called upon receiving GuildMessageCreateEvent, GuildMessageUpdateEvent or GuildMessageDeleteEvent
type GenericGuildMessageEvent struct {
	*GenericEvent
	MessageID snowflake.Snowflake
	Message   discord.Message
	ChannelID snowflake.Snowflake
	GuildID   snowflake.Snowflake
}

// Guild returns the discord.Guild the GenericGuildMessageEvent happened in.
// This will only check cached guilds!
func (e GenericGuildMessageEvent) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}

// Channel returns the discord.DMChannel where the GenericGuildMessageEvent happened
func (e GenericGuildMessageEvent) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}

// GuildMessageCreateEvent is called upon receiving a discord.Message in a Channel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing a discord.Message in a Channel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage discord.Message
}

// GuildMessageDeleteEvent is called upon deleting a discord.Message in a Channel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
