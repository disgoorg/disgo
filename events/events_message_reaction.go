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

// MessageReactionAddEvent indicates that an core.User added an core.MessageReaction to an core.Message in an core.GetChannel(this+++ requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEvent indicates that an core.User removed an core.MessageReaction from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvent
}

// MessageReactionRemoveEmojiEvent indicates someone removed all core.MessageReaction of a specific core.Emoji from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	Emoji discord.ReactionEmoji
}

// MessageReactionRemoveAllEvent indicates someone removed all core.MessageReaction(s) from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
