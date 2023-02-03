package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuildScheduledEvent is the base struct for all GuildScheduledEvents events.
type GenericGuildScheduledEvent struct {
	*GenericEvent
	GuildScheduled discord.GuildScheduledEvent
}

// GuildScheduledEventCreate is dispatched when a guild scheduled event is created.
type GuildScheduledEventCreate struct {
	*GenericGuildScheduledEvent
}

// GuildScheduledEventUpdate is dispatched when a guild scheduled event is updated.
type GuildScheduledEventUpdate struct {
	*GenericGuildScheduledEvent
	OldGuildScheduled discord.GuildScheduledEvent
}

// GuildScheduledEventDelete is dispatched when a guild scheduled event is deleted.
type GuildScheduledEventDelete struct {
	*GenericGuildScheduledEvent
}

// GenericGuildScheduledEventUser is the base struct for all GuildScheduledEventUser events.
type GenericGuildScheduledEventUser struct {
	*GenericEvent
	GuildScheduledEventID snowflake.ID
	UserID                snowflake.ID
	GuildID               snowflake.ID
}

// GuildScheduledEvent returns the discord.GuildScheduledEvent the event is for.
func (e *GenericGuildScheduledEventUser) GuildScheduledEvent() (discord.GuildScheduledEvent, bool) {
	return e.Client().Caches().GuildScheduledEvent(e.GuildID, e.GuildScheduledEventID)
}

// Member returns the Member who was added/removed from the GuildScheduledEvent from the cache.
func (e *GenericGuildScheduledEventUser) Member() (discord.Member, bool) {
	return e.Client().Caches().Member(e.GuildID, e.UserID)
}

// GuildScheduledEventUserAdd is dispatched when a user is added to a discord.GuildScheduledEvent.
type GuildScheduledEventUserAdd struct {
	*GenericGuildScheduledEventUser
}

// GuildScheduledEventUserRemove is dispatched when a user is removed from a discord.GuildScheduledEvent.
type GuildScheduledEventUserRemove struct {
	*GenericGuildScheduledEventUser
}
