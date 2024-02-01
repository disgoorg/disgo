package events

import (
	"github.com/snekROmonoro/snowflake"

	"github.com/snekROmonoro/disgo/discord"
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
	return e.Client().Caches().Guild(e.GuildId)
}

// Channel returns the discord.GuildMessageChannel the webhook was updated in.
// This will only return cached channels!
func (e *WebhooksUpdate) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().GuildMessageChannel(e.ChannelID)
}
