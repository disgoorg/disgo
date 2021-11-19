package events

import (
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMemberEvent generic core.Member event
type GenericGuildMemberEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
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
	GuildID discord.Snowflake
	User    *core.User
}

// GuildMemberTypingEvent indicates that a core.Member started typing in a core.BaseGuildMessageChannel(requires discord.GatewayIntentGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	UserID    discord.Snowflake
	GuildID   discord.Snowflake
	Timestamp time.Time
}

// MessageChannel returns the core.BaseGuildMessageChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) MessageChannel() core.BaseGuildMessageChannel {
	if ch := e.Bot().Caches.ChannelCache().Get(e.ChannelID); ch != nil {
		return ch.(core.BaseGuildMessageChannel)
	}
	return nil
}
