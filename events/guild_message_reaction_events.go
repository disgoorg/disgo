package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericGuildMessageReactionEvent is called upon receiving GuildMessageReactionAddEvent or GuildMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericEvent
	UserID    snowflake.Snowflake
	ChannelID snowflake.Snowflake
	MessageID snowflake.Snowflake
	GuildID   snowflake.Snowflake
	Emoji     discord.ReactionEmoji
}

func (e *GenericGuildMessageReactionEvent) Member() (discord.Member, bool) {
	return e.Client().Caches().Members().Get(e.GuildID, e.UserID)
}

// GuildMessageReactionAddEvent indicates that a discord.Member added a discord.ReactionEmoji to a discord.Message in a discord.GuildMessageChannel(requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
	Member discord.Member
}

// GuildMessageReactionRemoveEvent indicates that a discord.Member removed a discord.MessageReaction from a discord.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific discord.Emoji from a discord.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	MessageID snowflake.Snowflake
	GuildID   snowflake.Snowflake
	Emoji     discord.ReactionEmoji
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a discord.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	MessageID snowflake.Snowflake
	GuildID   snowflake.Snowflake
}
