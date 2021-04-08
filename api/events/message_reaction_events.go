package events

import "github.com/DisgoOrg/disgo/api"

type GenericReactionEvents struct {
	GenericMessageEvent
	UserID          api.Snowflake
	User            *api.User
	MessageReaction api.MessageReaction
}

type MessageReactionAddEvent struct {
	GenericReactionEvents
}

type MessageReactionRemoveEvent struct {
	GenericReactionEvents
}

type MessageReactionRemoveEmoteEvent struct {
	GenericMessageEvent
	MessageReaction api.MessageReaction
}

type MessageReactionRemoveAllEvent struct {
	GenericMessageEvent
}
