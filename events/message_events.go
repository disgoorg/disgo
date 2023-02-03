package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericMessage generic discord.Message event
type GenericMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   discord.Message
	ChannelID snowflake.ID
	GuildID   *snowflake.ID
}

// Channel returns the discord.GuildMessageChannel where the GenericMessage happened
func (e *GenericMessage) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().GuildMessageChannel(e.ChannelID)
}

// Guild returns the discord.Guild where the GenericMessage happened or nil if it happened in DMs
func (e *GenericMessage) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Client().Caches().Guild(*e.GuildID)
}

// MessageCreate indicates that a discord.Message got received
type MessageCreate struct {
	*GenericMessage
}

// MessageUpdate indicates that a discord.Message got update
type MessageUpdate struct {
	*GenericMessage
	OldMessage discord.Message
}

// MessageDelete indicates that a discord.Message got deleted
type MessageDelete struct {
	*GenericMessage
}
