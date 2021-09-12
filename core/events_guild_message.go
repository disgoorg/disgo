package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageEvent is called upon receiving GuildMessageCreateEvent, GuildMessageUpdateEvent or GuildMessageDeleteEvent
type GenericGuildMessageEvent struct {
	*GenericMessageEvent
	GuildID discord.Snowflake
}

// Guild returns the core.Guild the GenericGuildMessageEvent happened in
func (e GenericGuildMessageEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildMessageCreateEvent is called upon receiving an core.Message in an core.DMChannel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing an core.Message in an core.DMChannel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage *Message
}

// GuildMessageDeleteEvent is called upon deleting an core.Message in an core.DMChannel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
