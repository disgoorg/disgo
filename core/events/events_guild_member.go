package events

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericGuildMemberEvent generic core.Member event
type GenericGuildMemberEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	Member  discord.Member
}

// GuildMemberJoinEvent indicates that a core.Member joined the core.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that a core.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember discord.Member
}

// GuildMemberLeaveEvent indicates that a core.Member left the core.Guild
type GuildMemberLeaveEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	User    discord.User
	Member  discord.Member
}

// GuildMemberTypingStartEvent indicates that a core.Member started typing in a core.BaseGuildMessageChannel(requires discord.GatewayIntentGuildMessageTyping)
type GuildMemberTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	UserID    snowflake.Snowflake
	GuildID   snowflake.Snowflake
	Timestamp time.Time
	Member    discord.Member
}

// Channel returns the core.BaseGuildMessageChannel the GuildMemberTypingStartEvent happened in
func (e GuildMemberTypingStartEvent) Channel() (discord.BaseGuildMessageChannel, bool) {
	return e.Bot().Caches().Channels().GetBaseGuildMessageChannel(e.ChannelID)
}
