package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent(requires the core.GatewayIntentsDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID          discord.Snowflake
	User            User
	MessageReaction discord.MessageReaction
}

// DMMessageReactionAddEvent indicates that an core.User added an core.MessageReaction to an core.Message in an core.DMChannel(requires the core.GatewayIntentsDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEvent indicates that an core.User removed an core.MessageReaction from an core.Message in an core.DMChannel(requires the core.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all core.MessageReaction of a specific core.Emoji from an core.Message in an core.DMChannel(requires the core.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericDMMessageEvent
	MessageReaction discord.MessageReaction
}

// DMMessageReactionRemoveAllEvent indicates someone removed all core.MessageReaction(s) from an core.Message in an core.DMChannel(requires the core.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericDMMessageEvent
}
