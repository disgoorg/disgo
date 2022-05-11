package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericReactionEvent is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvent struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
	Emoji     discord.ReactionEmoji
}

// MessageReactionAddEvent indicates that a discord.User added a discord.MessageReaction to a discord.Message in a discord.Channel(this+++ requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvent
	Member *discord.Member
}

// MessageReactionRemoveEvent indicates that a discord.User removed a discord.MessageReaction from a discord.Message in a discord.GetChannel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific discord.Emoji from a discord.Message in a discord.Channel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
	Emoji     discord.ReactionEmoji
}

// MessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a discord.Message in a discord.Channel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
}
