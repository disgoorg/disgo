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

// Guild returns the api.Guild the GenericGuildMessageEvent happened in
func (e GenericGuildMessageEvent) Guild() *core.Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

// GuildMessageCreateEvent is called upon receiving an api.Message in an api.DMChannel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing an api.Message in an api.DMChannel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage *core.Message
}

// GuildMessageDeleteEvent is called upon deleting an api.Message in an api.DMChannel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
