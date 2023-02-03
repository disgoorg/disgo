package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuildChannel is called upon receiving GuildChannelCreate , GuildChannelUpdate or GuildChannelDelete
type GenericGuildChannel struct {
	*GenericEvent
	ChannelID snowflake.ID
	Channel   discord.GuildChannel
	GuildID   snowflake.ID
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e *GenericGuildChannel) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guild(e.GuildID)
}

// GuildChannelCreate indicates that a new Channel got created in a discord.Guild
type GuildChannelCreate struct {
	*GenericGuildChannel
}

// GuildChannelUpdate indicates that a Channel got updated in a discord.Guild
type GuildChannelUpdate struct {
	*GenericGuildChannel
	OldChannel discord.GuildChannel
}

// GuildChannelDelete indicates that a Channel got deleted in a discord.Guild
type GuildChannelDelete struct {
	*GenericGuildChannel
}

// GuildChannelPinsUpdate indicates a discord.Message got pinned or unpinned in a discord.GuildMessageChannel
type GuildChannelPinsUpdate struct {
	*GenericEvent
	GuildID             snowflake.ID
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
	OldLastPinTimestamp *time.Time
}
