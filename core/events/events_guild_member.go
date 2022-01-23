package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

// GenericGuildMemberEvent generic core.Member event
type GenericGuildMemberEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	Member  *core.Member
}

// GuildMemberJoinEvent indicates that a core.Member joined the core.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that a core.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *core.Member
}

// GuildMemberLeaveEvent indicates that a core.Member left the core.Guild
type GuildMemberLeaveEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	User    *core.User
	Member  *core.Member
}

// GuildMemberTypingStartEvent indicates that a core.Member started typing in a core.BaseGuildMessageChannel(requires discord.GatewayIntentGuildMessageTyping)
type GuildMemberTypingStartEvent struct {
	*GenericEvent
	ChannelID snowflake.Snowflake
	UserID    snowflake.Snowflake
	GuildID   snowflake.Snowflake
	Timestamp time.Time
	Member    *core.Member
}

// Channel returns the core.BaseGuildMessageChannel the GuildMemberTypingStartEvent happened in
func (e GuildMemberTypingStartEvent) Channel() core.BaseGuildMessageChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.BaseGuildMessageChannel)
	}
	return nil
}
