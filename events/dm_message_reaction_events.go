package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent (requires the discord.GatewayIntentDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionAddEvent indicates that a discord.User added a discord.MessageReaction to a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEvent indicates that a discord.User removed a discord.MessageReaction from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction(s) of a specific discord.Emoji from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a discord.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
}
