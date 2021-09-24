package core

import "github.com/DisgoOrg/disgo/discord"

type GuildIntegrationsUpdateEvent struct {
	*GenericEvent
	GuildId discord.Snowflake
}

func (e *GuildIntegrationsUpdateEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildId)
}
