package core

import "github.com/DisgoOrg/disgo/discord"

type WebhooksUpdateEvent struct {
	*GenericEvent
	GuildId   discord.Snowflake
	ChannelID discord.Snowflake
}

func (e *WebhooksUpdateEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildId)
}

func (e *WebhooksUpdateEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.GuildId)
}
