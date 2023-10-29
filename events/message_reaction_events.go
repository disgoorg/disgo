package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericReaction is called upon receiving MessageReactionAdd or MessageReactionRemove
type GenericReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	GuildID     *snowflake.ID
	Emoji       discord.PartialEmoji
	BurstColors []string
	Burst       bool
}

// MessageReactionAdd indicates that a discord.User added a discord.MessageReaction to a discord.Message in a discord.Channel(this+++ requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionAdd struct {
	*GenericReaction
	Member *discord.Member
}

// MessageReactionRemove indicates that a discord.User removed a discord.MessageReaction from a discord.Message in a discord.GetChannel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemove struct {
	*GenericReaction
}

// MessageReactionRemoveEmoji indicates someone removed all discord.MessageReaction of a specific discord.Emoji from a discord.Message in a discord.Channel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
	Emoji     discord.PartialEmoji
}

// MessageReactionRemoveAll indicates someone removed all discord.MessageReaction(s) from a discord.Message in a discord.Channel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
}
