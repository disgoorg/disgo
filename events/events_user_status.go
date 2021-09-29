package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// UserStatusUpdateEvent generic Status event
type UserStatusUpdateEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

// User returns the User that changed their Status.
// This will only check cached users!
func (g *UserStatusUpdateEvent) User() *core.User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}

// UserClientStatusUpdateEvent generic client-specific Status event
type UserClientStatusUpdateEvent struct {
	*GenericEvent
	UserID          discord.Snowflake
	OldClientStatus *discord.ClientStatus
	ClientStatus    discord.ClientStatus
}

// User returns the User that changed their Status.
// This will only check cached users!
func (g *UserClientStatusUpdateEvent) User() *core.User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}
