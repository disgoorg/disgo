package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	Channel   core.GuildChannel
	GuildID   discord.Snowflake
}

// Guild returns the core.Guild the event happened in.
// This will only check cached guilds!
func (e GenericGuildChannelEvent) Guild() *core.Guild {
	return e.Bot().Caches().Guilds().Get(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new Channel got created in a core.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that a Channel got updated in a core.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldChannel core.GuildChannel
}

// GuildChannelDeleteEvent indicates that a Channel got deleted in a core.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}

type GuildChannelPinsUpdateEvent struct {
	*GenericEvent
	GuildID             discord.Snowflake
	ChannelID           discord.Snowflake
	NewLastPinTimestamp *discord.Time
	OldLastPinTimestamp *discord.Time
}
