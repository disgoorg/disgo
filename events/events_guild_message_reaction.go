package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID discord.Snowflake
	Emoji  discord.ReactionEmoji
}

func (e *GenericGuildMessageReactionEvent) User() *core.User {
	return e.Bot().Caches.UserCache().Get(e.UserID)
}

func (e *GenericGuildMessageReactionEvent) Member() *core.Member {
	return e.Bot().Caches.MemberCache().Get(e.GuildID, e.UserID)
}

// GuildMessageReactionAddEvent indicates that an core.Member added an core.MessageReaction to an core.Message in an core.TextChannel(requires the core.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
	Member *core.Member
}

// GuildMessageReactionRemoveEvent indicates that an core.Member removed an core.MessageReaction from an core.Message in an core.TextChannel(requires the core.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all core.MessageReaction of a specific core.Emoji from an core.Message in an core.TextChannel(requires the core.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericGuildMessageEvent
	Emoji discord.ReactionEmoji
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all core.MessageReaction(s) from an core.Message in an core.TextChannel(requires the core.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericGuildMessageEvent
}
