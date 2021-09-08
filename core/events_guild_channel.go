package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericChannelEvent
	GuildID discord.Snowflake
}

// Guild returns the cached api.Guild the event happened in
func (e GenericGuildChannelEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new api.GetGuildChannel got created in an api.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that an api.GetGuildChannel got updated in an api.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldChannel *Channel
}

type GuildChannelPinsUpdateEvent struct {
	*GenericGuildChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// GuildChannelDeleteEvent indicates that an api.GetGuildChannel got deleted in an api.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}
