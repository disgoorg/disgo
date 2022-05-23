package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// WebhooksUpdate indicates that a guilds webhooks were updated.
type WebhooksUpdate struct {
	*GenericEvent
	GuildId   snowflake.ID
	ChannelID snowflake.ID
}

// Guild returns the Guild the webhook was updated in.
// This will only return cached guilds!
func (e *WebhooksUpdate) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildId)
}

// Channel returns the Channel the webhook was updated in.
// This will only return cached channels!
func (e *WebhooksUpdate) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}
