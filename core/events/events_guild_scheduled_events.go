package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericGuildScheduledEventEvent struct {
	*GenericEvent
	GuildScheduledEvent *core.GuildScheduledEvent
}

type GuildScheduledEventCreateEvent struct {
	*GenericGuildScheduledEventEvent
}

type GuildScheduledEventUpdateEvent struct {
	*GenericGuildScheduledEventEvent
	OldGuildScheduledEvent *core.GuildScheduledEvent
}

type GuildScheduledEventDeleteEvent struct {
	*GenericGuildScheduledEventEvent
}

type GenericGuildScheduledEventUserEvent struct {
	*GenericEvent
	GuildScheduledEventID discord.Snowflake
	UserID                discord.Snowflake
	GuildID               discord.Snowflake
}

func (e *GenericGuildScheduledEventUserEvent) GuildScheduledEvent() *core.GuildScheduledEvent {
	return e.bot.Caches.GuildScheduledEvents().Get(e.GuildScheduledEventID)
}

func (e *GenericGuildScheduledEventUserEvent) User() *core.User {
	return e.bot.Caches.Users().Get(e.UserID)
}

func (e *GenericGuildScheduledEventUserEvent) Member() *core.Member {
	return e.bot.Caches.Members().Get(e.GuildID, e.UserID)
}

type GuildScheduledEventUserAddEvent struct {
	*GenericGuildScheduledEventUserEvent
}

type GuildScheduledEventUserRemoveEvent struct {
	*GenericGuildScheduledEventUserEvent
}
