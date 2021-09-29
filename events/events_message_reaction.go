package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericReactionEvent is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvent struct {
	*GenericMessageEvent
	UserID discord.Snowflake
	Emoji  discord.ReactionEmoji
}

func (e *GenericReactionEvent) User() *core.User {
	return e.Bot().Caches.UserCache().Get(e.UserID)
}

// MessageReactionAddEvent indicates that a core.User added a discord.MessageReaction to a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEvent indicates that a core.User removed a discord.MessageReaction from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction of a specific core.Emoji from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	Emoji discord.ReactionEmoji
}

// MessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a core.Message in a Channel (requires the discord.GatewayIntentGuildMessageReactions and/or discord.GatewayIntentDirectMessageReactions)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
