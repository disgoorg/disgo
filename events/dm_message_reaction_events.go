package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMMessageReaction is called upon receiving DMMessageReactionAdd or DMMessageReactionRemove (requires the discord.GatewayIntentDirectMessageReactions)
type GenericDMMessageReaction struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionAdd indicates that a discord.User added a discord.MessageReaction to a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionAdd struct {
	*GenericDMMessageReaction
}

// DMMessageReactionRemove indicates that a discord.User removed a discord.MessageReaction from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemove struct {
	*GenericDMMessageReaction
}

// DMMessageReactionRemoveEmoji indicates someone removed all discord.MessageReaction(s) of a specific discord.Emoji from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionRemoveAll indicates someone removed all discord.MessageReaction(s) from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
}
