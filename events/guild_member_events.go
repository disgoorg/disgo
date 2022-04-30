package events

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildMemberEvent generic discord.Member event
type GenericGuildMemberEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	Member  discord.Member
}

// GuildMemberJoinEvent indicates that a discord.Member joined the discord.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that a discord.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember discord.Member
}

// GuildMemberLeaveEvent indicates that a discord.Member left the discord.Guild
type GuildMemberLeaveEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
	Member  discord.Member
}

// GuildMemberTypingStartEvent indicates that a discord.Member started typing in a discord.BaseGuildMessageChannel(requires discord.GatewayIntentGuildMessageTyping)
type GuildMemberTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	GuildID   snowflake.ID
	Timestamp time.Time
	Member    discord.Member
}

// Channel returns the discord.BaseGuildMessageChannel the GuildMemberTypingStartEvent happened in
func (e GuildMemberTypingStartEvent) Channel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID)
}
