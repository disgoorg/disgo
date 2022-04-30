package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	Channel   discord.GuildChannel
	GuildID   snowflake.ID
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e GenericGuildChannelEvent) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new Channel got created in a discord.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that a Channel got updated in a discord.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldChannel discord.GuildChannel
}

// GuildChannelDeleteEvent indicates that a Channel got deleted in a discord.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}

type GuildChannelPinsUpdateEvent struct {
	*GenericEvent
	GuildID             snowflake.ID
	ChannelID           snowflake.ID
	NewLastPinTimestamp *discord.Time
	OldLastPinTimestamp *discord.Time
}
