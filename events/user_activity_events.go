package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericUserActivityEvent generic Activity event
type GenericUserActivityEvent struct {
	*GenericEvent
	UserID   snowflake.ID
	GuildID  snowflake.ID
	Activity discord.Activity
}

// Member returns the Member that changed their Activity.
// This will only check cached members!
func (g *GenericUserActivityEvent) Member() (discord.Member, bool) {
	return g.Client().Caches().Members().Get(g.GuildID, g.UserID)
}

// Guild returns the Guild that changed their Activity.
// This will only check cached guilds!
func (g *GenericUserActivityEvent) Guild() (discord.Guild, bool) {
	return g.Client().Caches().Guilds().Get(g.UserID)
}

// UserActivityStartEvent indicates that a User started an Activity
type UserActivityStartEvent struct {
	*GenericUserActivityEvent
}

// UserActivityUpdateEvent indicates that a User updated their Activity
type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	OldActivity discord.Activity
}

// UserActivityStopEvent indicates that a User stopped an Activity
type UserActivityStopEvent struct {
	*GenericUserActivityEvent
}
