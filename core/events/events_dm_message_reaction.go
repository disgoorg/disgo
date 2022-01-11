package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent (requires the discord.GatewayIntentDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericEvent
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	Emoji     discord.ReactionEmoji
}

// User returns the User who owns the discord.MessageReaction.
// This will only check cached users!
func (e *GenericDMMessageReactionEvent) User() *core.User {
	return e.Bot().Caches().Users().Get(e.UserID)
}

// DMMessageReactionAddEvent indicates that a core.User added a discord.MessageReaction to a core.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEvent indicates that a core.User removed a discord.MessageReaction from a core.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all discord.MessageReaction(s) of a specific core.Emoji from a core.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	Emoji     discord.ReactionEmoji
}

// DMMessageReactionRemoveAllEvent indicates someone removed all discord.MessageReaction(s) from a core.Message in a Channel (requires the discord.GatewayIntentDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
}
