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

// UserUpdateEvent indicates that a core.User has updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *User
}

// UserTypingEvent indicates that a core.User started typing in a Channel (requires the discord.GatewayIntentDirectMessageTyping and/or discord.GatewayIntentGuildMessageTyping)
type UserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// Channel returns the Channel the core.User started typing in.
// This will only check cached channels!
func (e *UserTypingEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
