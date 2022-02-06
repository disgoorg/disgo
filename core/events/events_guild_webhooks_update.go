package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

type WebhooksUpdateEvent struct {
	*GenericEvent
	GuildId   snowflake.Snowflake
	ChannelID snowflake.Snowflake
}

// Guild returns the Guild the webhook was updated in.
// This will only check cached guilds!
func (e *WebhooksUpdateEvent) Guild() *core.Guild {
	return e.Bot().Caches.Guilds().Get(e.GuildId)
}

func (e *WebhooksUpdateEvent) Channel() core.GuildMessageChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.GuildMessageChannel)
	}
	return nil
}