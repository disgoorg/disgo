package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageEvent is called upon receiving GuildMessageCreateEvent, GuildMessageUpdateEvent or GuildMessageDeleteEvent
type GenericGuildMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *core.Message
	ChannelID discord.Snowflake
	GuildID   discord.Snowflake
}

// Guild returns the core.Guild the GenericGuildMessageEvent happened in
func (e GenericGuildMessageEvent) Guild() *core.Guild {
	return e.Bot().Caches.Guilds().Get(e.GuildID)
}

// Channel returns the core.DMChannel where the GenericGuildMessageEvent happened
func (e GenericGuildMessageEvent) Channel() core.GuildMessageChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.GuildMessageChannel)
	}
	return nil
}

// GuildMessageCreateEvent is called upon receiving an core.Message in an core.DMChannel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing an core.Message in an core.DMChannel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage *core.Message
}

// GuildMessageDeleteEvent is called upon deleting an core.Message in an core.DMChannel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
