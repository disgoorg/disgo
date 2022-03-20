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

// UserUpdateEvent  indicates that a discord.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser discord.User
}

// UserTypingStartEvent indicates that a discord.User started typing in a discord.DMChannel or discord.MessageChanel(requires the discord.GatewayIntentDirectMessageTyping and/or discord.GatewayIntentGuildMessageTyping)
type UserTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	GuildID   *snowflake.Snowflake
	UserID    snowflake.Snowflake
	Timestamp time.Time
}

// Channel returns the discord.MessageChannel the discord.User started typing in
func (e *UserTypingStartEvent) Channel() (discord.MessageChannel, bool) {
	return e.Bot().Caches().Channels().GetMessageChannel(e.ChannelID)
}
