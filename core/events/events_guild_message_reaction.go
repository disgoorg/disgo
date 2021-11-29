package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageReactionEvent is called upon receiving GuildMessageReactionAddEvent or GuildMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   discord.Snowflake
	Emoji     discord.ReactionEmoji
}

func (e *GenericGuildMessageReactionEvent) User() *core.User {
	return e.Bot().Caches.Users().Get(e.UserID)
}

func (e *GenericGuildMessageReactionEvent) Member() *core.Member {
	return e.Bot().Caches.Members().Get(e.GuildID, e.UserID)
}

// GuildMessageReactionAddEvent indicates that a core.Member added a discord.ReactionEmoji to a core.Message in a core.GuildMessageChannel(requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
	Member *core.Member
}

// GuildMessageReactionRemoveEvent indicates that a core.Member removed a discord.MessageReaction from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific core.Emoji from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   discord.Snowflake
	Emoji     discord.ReactionEmoji
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   discord.Snowflake
}
