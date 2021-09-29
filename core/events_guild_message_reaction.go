package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageReactionEvent is called upon receiving GuildMessageReactionAddEvent or GuildMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID discord.Snowflake
	Member *Member
	Emoji  discord.ReactionEmoji
}

// GuildMessageReactionAddEvent indicates that a core.Member added a discord.MessageReaction to a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEvent indicates that a core.Member removed a discord.MessageReaction from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific core.Emoji from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericGuildMessageEvent
	Emoji discord.ReactionEmoji
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericGuildMessageEvent
}
