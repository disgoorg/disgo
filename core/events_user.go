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

// UserUpdateEvent  indicates that an core.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *User
}

// UserTypingEvent indicates that an core.User started typing in an core.DMChannel or core.TextChannel(requires the core.GatewayIntentsDirectMessageTyping and/or core.GatewayIntentsGuildMessageTyping)
type UserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// Channel returns the core.GetChannel the core.User started typing in
func (e *UserTypingEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
