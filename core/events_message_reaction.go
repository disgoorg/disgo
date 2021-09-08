package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericReactionEvents is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvents struct {
	*GenericMessageEvent
	UserID          discord.Snowflake
	User            *User
	MessageReaction discord.MessageReaction
}

// MessageReactionAddEvent indicates that an core.User added an core.MessageReaction to an core.Message in an core.GetChannel(this+++ requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEvent indicates that an core.User removed an core.MessageReaction from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEmojiEvent indicates someone removed all core.MessageReaction of a specific core.Emoji from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	MessageReaction discord.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all core.MessageReaction(s) from an core.Message in an core.GetChannel(requires the core.GatewayIntentsGuildMessageReactions and/or core.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
