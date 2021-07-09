package events

import "github.com/DisgoOrg/disgo/api"

// GenericMessageReactionEvent is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericMessageReactionEvent struct {
	*GenericMessageEvent
	Emoji *api.Emoji
}

type GenericMessageUserReactionEvent struct {
	*GenericMessageReactionEvent
	UserID api.Snowflake
}

func (e *GenericMessageUserReactionEvent) User() *api.User {
	return e.Disgo().Cache().User(e.UserID)
}

// MessageReactionAddEvent indicates that a api.User added a api.MessageReaction to a api.Message in a api.Channel(this requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericMessageUserReactionEvent
}

// MessageReactionRemoveEvent indicates that a api.User removed a api.MessageReaction from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericMessageUserReactionEvent
}

// MessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageReactionEvent
	MessageReaction api.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
