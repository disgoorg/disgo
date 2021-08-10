package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID          discord.Snowflake
	Member          *core.Member
	MessageReaction discord.MessageReaction
}

// GuildMessageReactionAddEvent indicates that an api.Member added an api.MessageReaction to an api.Message in an api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEvent indicates that an api.Member removed an api.MessageReaction from an api.Message in an api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from an api.Message in an api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericGuildMessageEvent
	MessageReaction discord.MessageReaction
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from an api.Message in an api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericGuildMessageEvent
}
