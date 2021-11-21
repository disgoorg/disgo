package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericUserActivityEvent struct {
	*GenericEvent
	UserID   discord.Snowflake
	GuildID  discord.Snowflake
	Activity discord.Activity
}

func (g *GenericUserActivityEvent) User() *core.User {
	return g.Bot().Caches.UserCache().Get(g.UserID)
}

func (g *GenericUserActivityEvent) Member() *core.Member {
	return g.Bot().Caches.Members().Get(g.GuildID, g.UserID)
}

func (g *GenericUserActivityEvent) Guild() *core.Guild {
	return g.Bot().Caches.Guilds().Get(g.UserID)
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
