package events

import "github.com/DisgoOrg/disgo/api"

// GenericDMMessageReactionEvent is called upon receiving DMMessageReactionAddEvent or DMMessageReactionRemoveEvent(requires the api.GatewayIntentsDirectMessageReactions)
type GenericDMMessageReactionEvent struct {
	*GenericDMMessageEvent
	Emoji *api.Emoji
}

type GenericDMMessageUserReactionEvent struct {
	*GenericDMMessageReactionEvent
	UserID api.Snowflake
}

func (e *GenericDMMessageUserReactionEvent) User() *api.User {
	return e.Disgo().Cache().User(e.UserID)
}

// DMMessageReactionAddEvent indicates that a api.User added a api.MessageReaction to a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionAddEvent struct {
	*GenericDMMessageUserReactionEvent
}

// DMMessageReactionRemoveEvent indicates that a api.User removed a api.MessageReaction from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEvent struct {
	*GenericDMMessageUserReactionEvent
}

// DMMessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveEmojiEvent struct {
	*GenericDMMessageReactionEvent
}

// DMMessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.DMChannel(requires the api.GatewayIntentsDirectMessageReactions)
type DMMessageReactionRemoveAllEvent struct {
	*GenericDMMessageEvent
}
