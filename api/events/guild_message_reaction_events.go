package events

import "github.com/DisgoOrg/disgo/api"

// GenericGuildMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent
type GenericGuildMessageReactionEvent struct {
	*GenericGuildMessageEvent
	Emoji *api.Emoji
}

type GenericGuildMessageUserReactionEvent struct {
	*GenericGuildMessageReactionEvent
	UserID api.Snowflake
}

func (e *GenericGuildMessageUserReactionEvent) User() *api.User {
	return e.Disgo().Cache().User(e.UserID)
}

// GuildMessageReactionAddEvent indicates that a api.Member added a api.MessageReaction to a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionAddEvent struct {
	*GenericGuildMessageUserReactionEvent
	Member *api.Member
}

// GuildMessageReactionRemoveEvent indicates that a api.Member removed a api.MessageReaction from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEvent struct {
	*GenericGuildMessageUserReactionEvent
}

// GuildMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveEmojiEvent struct {
	*GenericGuildMessageReactionEvent
}

// GuildMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.TextChannel(requires the api.GatewayIntentsGuildMessageReactions)
type GuildMessageReactionRemoveAllEvent struct {
	*GenericGuildMessageEvent
}
