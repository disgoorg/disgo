package events

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
)

type WebhooksUpdateEvent struct {
	*GenericEvent
	GuildId   snowflake.Snowflake
	ChannelID snowflake.Snowflake
}

// Guild returns the Guild the webhook was updated in.
// This will only check cached guilds!
func (e *WebhooksUpdateEvent) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildId)
}

func (e *WebhooksUpdateEvent) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}
