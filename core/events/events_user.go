package events

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingStartEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID snowflake.Snowflake
	User   discord.User
}

// UserUpdateEvent  indicates that a core.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser discord.User
}

// UserTypingStartEvent indicates that a core.User started typing in a core.DMChannel or core.MessageChanel(requires the discord.GatewayIntentDirectMessageTyping and/or discord.GatewayIntentGuildMessageTyping)
type UserTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	GuildID   *snowflake.Snowflake
	UserID    snowflake.Snowflake
	Timestamp time.Time
}

// MessageChannel returns the core.GetChannel the core.User started typing in
func (e *UserTypingStartEvent) MessageChannel() (discord.MessageChannel, bool) {
	return e.Bot().Caches().Channels().GetMessageChannel(e.ChannelID)
}
