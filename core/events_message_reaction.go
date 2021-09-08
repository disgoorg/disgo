package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericReactionEvents is called upon receiving MessageReactionAddEvent or MessageReactionRemoveEvent
type GenericReactionEvents struct {
	*GenericMessageEvent
	UserID          discord.Snowflake
	User            *User
	MessageReaction discord.MessageReaction
}

// MessageReactionAddEvent indicates that an api.User added an api.MessageReaction to an api.Message in an api.GetChannel(this+++ requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionAddEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEvent indicates that an api.User removed an api.MessageReaction from an api.Message in an api.GetChannel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEvent struct {
	*GenericReactionEvents
}

// MessageReactionRemoveEmojiEvent indicates someone removed all api.MessageReaction of a specific api.Emoji from an api.Message in an api.GetChannel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactions)
type MessageReactionRemoveEmojiEvent struct {
	*GenericMessageEvent
	MessageReaction discord.MessageReaction
}

// MessageReactionRemoveAllEvent indicates someone removed all api.MessageReaction(s) from an api.Message in an api.GetChannel(requires the api.GatewayIntentsGuildMessageReactions and/or api.GatewayIntentsDirectMessageReactionss)
type MessageReactionRemoveAllEvent struct {
	*GenericMessageEvent
}
