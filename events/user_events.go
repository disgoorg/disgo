package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericUser is called upon receiving UserUpdate or UserTypingStart
type GenericUser struct {
	*GenericEvent
	UserID snowflake.ID
	User   discord.User
}

// UserUpdate  indicates that a discord.User updated
type UserUpdate struct {
	*GenericUser
	OldUser discord.User
}

// UserTypingStart indicates that a discord.User started typing in a discord.DMChannel or discord.MessageChanel(requires the gateway.IntentDirectMessageTyping and/or gateway.IntentGuildMessageTyping)
type UserTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	GuildID   *snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}

// Channel returns the discord.GuildMessageChannel the discord.User started typing in
func (e *UserTypingStart) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().GuildMessageChannel(e.ChannelID)
}
