package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID discord.Snowflake
	User   *User
}

// UserUpdateEvent  indicates that an api.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *User
}

// UserTypingEvent indicates that an api.User started typing in an api.DMChannel or api.TextChannel(requires the api.GatewayIntentsDirectMessageTyping and/or api.GatewayIntentsGuildMessageTyping)
type UserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// Channel returns the api.GetChannel the api.User started typing in
func (e *UserTypingEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
