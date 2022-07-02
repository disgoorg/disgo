package events

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
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
type GuildMemberTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	GuildID   snowflake.ID
	Timestamp time.Time
	Member    discord.Member
}

// Channel returns the discord.BaseGuildMessageChannel the GuildMemberTypingStart happened in
func (e GuildMemberTypingStart) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}
