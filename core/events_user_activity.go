package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type GenericUserActivityEvent struct {
	*GenericEvent
	UserID   discord.Snowflake
	GuildID  discord.Snowflake
	Activity discord.Activity
}

func (g *GenericUserActivityEvent) User() *User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}

func (g *GenericUserActivityEvent) Member() *Member {
	return g.Bot().Caches.MemberCache().Get(g.GuildID, g.UserID)
}

func (g *GenericUserActivityEvent) Guild() *Guild {
	return g.Bot().Caches.GuildCache().Get(g.UserID)
}

type UserActivityStartEvent struct {
	*GenericUserActivityEvent
}

type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	OldActivity discord.Activity
}

type UserActivityStopEvent struct {
	*GenericUserActivityEvent
}
