package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent(requires the core.GatewayIntentsDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	Emoji     discord.ReactionEmoji
}

func (e *GenericDMMessageReactionEvent) User() *core.User {
	return e.Bot().Caches.Users().Get(e.UserID)
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
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionRemoveAllEvent indicates someone removed all core.MessageReaction(s) from an core.Message in an core.DMChannel(requires the core.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
}
