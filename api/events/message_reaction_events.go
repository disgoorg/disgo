package events

import "github.com/DisgoOrg/disgo/api"

// GenericReactionEvents is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvents struct {
	*GenericMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

// MessageReactionAddEvent indicates that a api.User added a api.MessageReaction to a api.Message in a api.Channel(this+++ requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEvent indicates that a api.User removed a api.MessageReaction from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	MessageReaction api.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
