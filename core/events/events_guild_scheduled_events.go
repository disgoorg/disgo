package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type GenericGuildScheduledEventEvent struct {
	*GenericEvent
	GuildScheduledEvent discord.GuildScheduledEvent
}

type GuildScheduledEventCreateEvent struct {
	*GenericGuildScheduledEventEvent
}

type GuildScheduledEventUpdateEvent struct {
	*GenericGuildScheduledEventEvent
	OldGuildScheduledEvent discord.GuildScheduledEvent
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

func (e *GenericGuildScheduledEventUserEvent) GuildScheduledEvent() (discord.GuildScheduledEvent, bool) {
	return e.bot.Caches().GuildScheduledEvents().Get(e.GuildScheduledEventID)
}

func (e *GenericGuildScheduledEventUserEvent) User() (discord.User, bool) {
	return e.bot.Caches().Users().Get(e.UserID)
}

func (e *GenericGuildScheduledEventUserEvent) Member() (discord.Member, bool) {
	return e.bot.Caches().Members().Get(e.GuildID, e.UserID)
}

type GuildScheduledEventUserAddEvent struct {
	*GenericGuildScheduledEventUserEvent
}

type GuildScheduledEventUserRemoveEvent struct {
	*GenericGuildScheduledEventUserEvent
}
