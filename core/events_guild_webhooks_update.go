package core

import "github.com/DisgoOrg/disgo/discord"

type WebhooksUpdateEvent struct {
	*GenericEvent
	GuildId   discord.Snowflake
	ChannelID discord.Snowflake
}

// Guild returns the Guild the webhook was updated in.
// This will only check cached guilds!
func (e *WebhooksUpdateEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildId)
}

// Channel returns the Channel the webhook was created for.
// This will only check cached channels!
func (e *WebhooksUpdateEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.GuildId)
}
