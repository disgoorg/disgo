package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageEvent is called upon receiving GuildMessageCreateEvent, GuildMessageUpdateEvent or GuildMessageDeleteEvent
type GenericGuildMessageEvent struct {
	*GenericMessageEvent
	GuildID discord.Snowflake
}

// Guild returns the core.Guild the GenericGuildMessageEvent happened in.
// This will only check cached guilds!
func (e GenericGuildMessageEvent) Guild() *core.Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildMessageCreateEvent is called upon receiving a core.Message in a Channel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing a core.Message in a Channel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage *core.Message
}

// GuildMessageDeleteEvent is called upon deleting a core.Message in a Channel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
