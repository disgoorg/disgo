package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericReactionEvents is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvents struct {
	*GenericMessageEvent
	UserID          discord.Snowflake
	User            *core.User
	MessageReaction discord.MessageReaction
}

// MessageReactionAddEvent indicates that an api.User added an api.MessageReaction to an api.Message in an api.Channel(this+++ requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEvent indicates that an api.User removed an api.MessageReaction from an api.Message in an api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from an api.Message in an api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	MessageReaction discord.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from an api.Message in an api.Channel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
