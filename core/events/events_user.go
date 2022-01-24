package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingStartEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID snowflake.Snowflake
	User   *core.User
}

// UserUpdateEvent  indicates that a core.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *core.User
}

// UserTypingStartEvent indicates that a core.User started typing in a core.DMChannel or core.MessageChanel(requires the discord.GatewayIntentDirectMessageTyping and/or discord.GatewayIntentGuildMessageTyping)
type UserTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	GuildID   *snowflake.Snowflake
	UserID    snowflake.Snowflake
	Timestamp time.Time
}

// Channel returns the core.GetChannel the core.User started typing in
func (e *UserTypingStartEvent) Channel() core.MessageChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.MessageChannel)
	}
	return nil
}
