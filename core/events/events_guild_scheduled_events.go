package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
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
	GuildScheduledEventID snowflake.Snowflake
	UserID                snowflake.Snowflake
	GuildID               snowflake.Snowflake
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
