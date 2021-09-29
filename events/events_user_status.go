package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type UserStatusUpdateEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

func (g *UserStatusUpdateEvent) User() *core.User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}

type UserClientStatusUpdateEvent struct {
	*GenericEvent
	UserID          discord.Snowflake
	OldClientStatus *discord.ClientStatus
	ClientStatus    discord.ClientStatus
}

func (g *UserClientStatusUpdateEvent) User() *core.User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}
