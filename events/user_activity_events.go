package events

import (
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type PresenceUpdate struct {
	*GenericEvent
	gateway.EventPresenceUpdate
}

// GenericUserActivity generic Activity event
type GenericUserActivity struct {
	*GenericEvent
	UserID   snowflake.ID
	GuildID  snowflake.ID
	Activity discord.Activity
}

// Member returns the Member that changed their Activity.
// This will only check cached members!
func (g *GenericUserActivity) Member() (discord.Member, bool) {
	return g.Client().Caches().Member(g.GuildID, g.UserID)
}

// Guild returns the Guild that changed their Activity.
// This will only check cached guilds!
func (g *GenericUserActivity) Guild() (discord.Guild, bool) {
	return g.Client().Caches().Guild(g.UserID)
}

// UserActivityStart indicates that a User started an Activity
type UserActivityStart struct {
	*GenericUserActivity
}

// UserActivityUpdate indicates that a User updated their Activity
type UserActivityUpdate struct {
	*GenericUserActivity
	OldActivity discord.Activity
}

// UserActivityStop indicates that a User stopped an Activity
type UserActivityStop struct {
	*GenericUserActivity
}
