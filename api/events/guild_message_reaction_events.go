package events

import "github.com/DisgoOrg/disgo/api"

type GenericGuildMessageReactionEvent struct {
	GenericGuildMessageEvent
	UserID          api.Snowflake
	Member          *api.Member
	MessageReaction api.MessageReaction
}

type GuildMessageReactionAddEvent struct {
	GenericGuildMessageReactionEvent
}

type GuildMessageReactionRemoveEvent struct {
	GenericGuildMessageReactionEvent
}

type GuildMessageReactionRemoveEmoteEvent struct {
	GenericGuildMessageEvent
	MessageReaction api.MessageReaction
}

type GuildMessageReactionRemoveAllEvent struct {
	GenericGuildMessageEvent
}
