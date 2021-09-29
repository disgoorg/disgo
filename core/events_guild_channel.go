package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericChannelEvent
	GuildID discord.Snowflake
}

// Guild returns the core.Guild the event happened in.
// This will only check cached guilds!
func (e GenericGuildChannelEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new Channel got created in a core.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that a Channel got updated in a core.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldChannel *Channel
}

type GuildChannelPinsUpdateEvent struct {
	*GenericGuildChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// GuildChannelDeleteEvent indicates that a Channel got deleted in a core.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}
