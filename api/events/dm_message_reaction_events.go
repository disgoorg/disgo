package events

import "github.com/DisgoOrg/disgo/api"

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent(requires the api.GatewayIntentsDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

// DMMessageReactionAddEvent indicates that a api.User added a api.MessageReaction to a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEvent indicates that a api.User removed a api.MessageReaction from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericDMMessageEvent
	MessageReaction api.MessageReaction
}

// DMMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericDMMessageEvent
}
