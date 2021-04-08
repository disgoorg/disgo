package events

import "github.com/DisgoOrg/disgo/api"

type GenericDMMessageReactionEvent struct {
	GenericGuildMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

type DMMessageReactionAddEvent struct {
	GenericDMMessageReactionEvent
}

type DMMessageReactionRemoveEvent struct {
	GenericDMMessageReactionEvent
}

type DMMessageReactionRemoveEmoteEvent struct {
	GenericDMMessageEvent
	MessageReaction api.MessageReaction
}

type DMMessageReactionRemoveAllEvent struct {
	GenericDMMessageEvent
}
