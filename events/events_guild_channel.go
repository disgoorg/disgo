package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericChannelEvent
	GuildID discord.Snowflake
}

// Guild returns the cached core.Guild the event happened in
func (e GenericGuildChannelEvent) Guild() *core.Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new core.GetGuildChannel got created in an core.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that an core.GetGuildChannel got updated in an core.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldChannel *core.Channel
}

type GuildChannelPinsUpdateEvent struct {
	*GenericGuildChannelEvent
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}

// GuildChannelDeleteEvent indicates that an core.GetGuildChannel got deleted in an core.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}
