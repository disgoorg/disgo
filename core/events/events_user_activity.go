package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericUserActivityEvent generic Activity event
type GenericUserActivityEvent struct {
	*GenericEvent
	UserID   discord.Snowflake
	GuildID  discord.Snowflake
	Activity discord.Activity
}

// User returns the User that changed their Activity.
// This will only check cached users!
func (g *GenericUserActivityEvent) User() *core.User {
	return g.Bot().Caches.Users().Get(g.UserID)
}

// Member returns the Member that changed their Activity.
// This will only check cached members!
func (g *GenericUserActivityEvent) Member() *core.Member {
	return g.Bot().Caches.Members().Get(g.GuildID, g.UserID)
}

// Guild returns the Guild that changed their Activity.
// This will only check cached guilds!
func (g *GenericUserActivityEvent) Guild() *core.Guild {
	return g.Bot().Caches.Guilds().Get(g.UserID)
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
