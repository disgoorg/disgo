package core

import "github.com/DisgoOrg/disgo/discord"

type UserStatusUpdateEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

func (g *UserStatusUpdateEvent) User() *User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}

type UserClientStatusUpdateEvent struct {
	*GenericEvent
	UserID          discord.Snowflake
	OldClientStatus *discord.ClientStatus
	ClientStatus    discord.ClientStatus
}

func (g *UserClientStatusUpdateEvent) User() *User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}
