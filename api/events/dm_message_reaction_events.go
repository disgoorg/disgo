package events

import "github.com/DisgoOrg/disgo/api"

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent(requires the api.GatewayIntentsDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericGuildMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

// DMMessageReactionAddEvent indicates that an api.User added an api.MessageReaction to an api.Message in an api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEvent indicates that an api.User removed an api.MessageReaction from an api.Message in an api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from an api.Message in an api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericDMMessageEvent
	MessageReaction api.MessageReaction
}

// DMMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from an api.Message in an api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericDMMessageEvent
}
