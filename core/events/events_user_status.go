package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// UserStatusUpdateEvent generic Status event
type UserStatusUpdateEvent struct {
	*GenericEvent
	UserID    snowflake.Snowflake
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

// User returns the User that changed their Status.
// This will only check cached users!
func (g *UserStatusUpdateEvent) User() (discord.User, bool) {
	return g.Bot().Caches().Users().Get(g.UserID)
}

// UserClientStatusUpdateEvent generic client-specific Status event
type UserClientStatusUpdateEvent struct {
	*GenericEvent
	UserID          snowflake.Snowflake
	OldClientStatus *discord.ClientStatus
	ClientStatus    discord.ClientStatus
}

// User returns the User that changed their Status.
// This will only check cached users!
func (g *UserClientStatusUpdateEvent) User() (discord.User, bool) {
	return g.Bot().Caches().Users().Get(g.UserID)
}
