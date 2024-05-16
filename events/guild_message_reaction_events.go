package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuildMessageReaction is called upon receiving GuildMessageReactionAdd or GuildMessageReactionRemove
type GenericGuildMessageReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	GuildID     snowflake.ID
	Emoji       discord.PartialEmoji
	BurstColors []string
	Burst       bool
}

// Member returns the Member that reacted to the discord.Message from the cache.
func (e *GenericGuildMessageReaction) Member() (discord.Member, bool) {
	return e.Client().Caches().Member(e.GuildID, e.UserID)
}

// GuildMessageReactionAdd indicates that a discord.Member added a discord.PartialEmoji to a discord.Message in a discord.GuildMessageChannel(requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionAdd struct {
	*GenericGuildMessageReaction
	Member          discord.Member
	MessageAuthorID *snowflake.ID
}

// GuildMessageReactionRemove indicates that a discord.Member removed a discord.MessageReaction from a discord.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemove struct {
	*GenericGuildMessageReaction
}

// GuildMessageReactionRemoveEmoji indicates someone removed all discord.MessageReaction of a specific discord.Emoji from a discord.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   snowflake.ID
	Emoji     discord.PartialEmoji
}

// GuildMessageReactionRemoveAll indicates someone removed all discord.MessageReaction(s) from a discord.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   snowflake.ID
}
