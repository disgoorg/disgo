package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericReactionEvent is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   *discord.Snowflake
	Emoji     discord.ReactionEmoji
}

func (e *GenericReactionEvent) User() *core.User {
	return e.Bot().Caches().Users().Get(e.UserID)
}

// MessageReactionAddEvent indicates that a core.User added a discord.MessageReaction to a core.Message in a core.Channel(this+++ requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvent
	Member *core.Member
}

// MessageReactionRemoveEvent indicates that a core.User removed a discord.MessageReaction from a core.Message in a core.GetChannel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific core.Emoji from a core.Message in a core.Channel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   *discord.Snowflake
	Emoji     discord.ReactionEmoji
}

// MessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a core.Message in a core.Channel(requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   *discord.Snowflake
}
