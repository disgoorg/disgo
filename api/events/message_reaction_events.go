package events

import "github.com/DisgoOrg/disgo/api"

// GenericReactionEvents is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvents struct {
	GenericMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

// MessageReactionAddEvent indicates that a api.User added a api.MessageReaction to a api.Message in a api.Channel(this+++ requires the api.IntentsGuildMessageReactions and/or api.IntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	GenericReactionEvents
}

// MessageReactionRemoveEvent indicates that a api.User removed a api.MessageReaction from a api.Message in a api.Channel(requires the api.IntentsGuildMessageReactions and/or api.IntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	GenericReactionEvents
}

// MessageReactionRemoveEmoteEvent indicates someone removed all api.MessageReaction of a specific api.Emote from a api.Message in a api.Channel(requires the api.IntentsGuildMessageReactions and/or api.IntentsDirectMessageReactions)
type MessageReactionRemoveEmoteEvent struct {
	GenericMessageEvent
	MessageReaction api.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from a api.Message in a api.Channel(requires the api.IntentsGuildMessageReactions and/or api.IntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	GenericMessageEvent
}
