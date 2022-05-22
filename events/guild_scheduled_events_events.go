package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GenericGuildScheduledEvent struct {
	*GenericEvent
	GuildScheduled discord.GuildScheduledEvent
}

type GuildScheduledEventCreate struct {
	*GenericGuildScheduledEvent
}
type GuildScheduledEventUpdate struct {
	*GenericGuildScheduledEvent
	OldGuildScheduled discord.GuildScheduledEvent
}
type GuildScheduledEventDelete struct {
	*GenericGuildScheduledEvent
}
type GenericGuildScheduledEventUser struct {
	*GenericEvent
	GuildScheduledEventID snowflake.ID
	UserID                snowflake.ID
	GuildID               snowflake.ID
}

func (e *GenericGuildScheduledEventUser) GuildScheduledEvent() (discord.GuildScheduledEvent, bool) {
	return e.Client().Caches().GuildScheduledEvents().Get(e.GuildID, e.GuildScheduledEventID)
}

func (e *GenericGuildScheduledEventUser) Member() (discord.Member, bool) {
	return e.Client().Caches().Members().Get(e.GuildID, e.UserID)
}

type GuildScheduledEventUserAdd struct {
	*GenericGuildScheduledEventUser
}
type GuildScheduledEventUserRemove struct {
	*GenericGuildScheduledEventUser
}
