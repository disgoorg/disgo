package events

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingStartEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID snowflake.ID
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
	ChannelID snowflake.ID
	GuildID   *snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}

// Channel returns the discord.MessageChannel the discord.User started typing in
func (e *UserTypingStartEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID)
}
