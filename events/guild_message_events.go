package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildMessage is called upon receiving GuildMessageCreate , GuildMessageUpdate or GuildMessageDelete
type GenericGuildMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   discord.Message
	ChannelID snowflake.ID
	GuildID   snowflake.ID
}

// Guild returns the discord.Guild the GenericGuildMessage happened in.
// This will only check cached guilds!
func (e GenericGuildMessage) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}

// Channel returns the discord.DMChannel where the GenericGuildMessage happened
func (e GenericGuildMessage) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}

// GuildMessageCreate is called upon receiving a discord.Message in a Channel
type GuildMessageCreate struct {
	*GenericGuildMessage
}

// GuildMessageUpdate is called upon editing a discord.Message in a Channel
type GuildMessageUpdate struct {
	*GenericGuildMessage
	OldMessage discord.Message
}

// GuildMessageDelete is called upon deleting a discord.Message in a Channel
type GuildMessageDelete struct {
	*GenericGuildMessage
}
