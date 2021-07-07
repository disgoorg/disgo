package events

import "github.com/DisgoOrg/disgo/api"

// GenericGuildMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID          api.Snowflake
	Member          *api.Member
	MessageReaction api.MessageReaction
}

// GuildMessageReactionAddEvent indicates that a api.Member added a api.MessageReaction to a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEvent indicates that a api.Member removed a api.MessageReaction from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericGuildMessageEvent
	MessageReaction api.MessageReaction
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericGuildMessageEvent
}
