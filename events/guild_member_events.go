package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuildMember generic discord.Member event
type GenericGuildMember struct {
	*GenericEvent
	GuildID snowflake.ID
	Member  discord.Member
}

// GuildMemberJoin indicates that a discord.Member joined the discord.Guild
type GuildMemberJoin struct {
	*GenericGuildMember
}

// GuildMemberUpdate indicates that a discord.Member updated
type GuildMemberUpdate struct {
	*GenericGuildMember
	OldMember discord.Member
}

// GuildMemberLeave indicates that a discord.Member left the discord.Guild
type GuildMemberLeave struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
	Member  discord.Member
}

// GuildMemberTypingStart indicates that a discord.Member started typing in a discord.BaseGuildMessageChannel(requires gateway.IntentGuildMessageTyping)
// Member will be empty when event is triggered by [Clyde bot]
//
// [Clyde bot]: https://support.discord.com/hc/en-us/articles/13066317497239-Clyde-Discord-s-AI-Chatbot
type GuildMemberTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	GuildID   snowflake.ID
	Timestamp time.Time
	Member    discord.Member
}

// Channel returns the discord.GuildMessageChannel the GuildMemberTypingStart happened in
func (e *GuildMemberTypingStart) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().GuildMessageChannel(e.ChannelID)
}
