package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericDMMessageReaction is called upon receiving DMMessageReactionAdd or DMMessageReactionRemove (requires the gateway.IntentDirectMessageReactions)
type GenericDMMessageReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	Emoji       discord.PartialEmoji
	BurstColors []string
	Burst       bool
}

// DMMessageReactionAdd indicates that a discord.User added a discord.MessageReaction to a discord.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionAdd struct {
	*GenericDMMessageReaction
	MessageAuthorID *snowflake.ID
}

// DMMessageReactionRemove indicates that a discord.User removed a discord.MessageReaction from a discord.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemove struct {
	*GenericDMMessageReaction
}

// DMMessageReactionRemoveEmoji indicates someone removed all discord.MessageReaction(s) of a specific discord.Emoji from a discord.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     discord.PartialEmoji
}

// DMMessageReactionRemoveAll indicates someone removed all discord.MessageReaction(s) from a discord.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
}
