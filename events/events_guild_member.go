package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMemberEvent generic core.Member event
type GenericGuildMemberEvent struct {
	*GenericGuildEvent
	Member *core.Member
}

// GuildMemberJoinEvent indicates that a core.Member joined a core.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that a core.Member has updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *core.Member
}

// GuildMemberLeaveEvent indicates that a core.Member left a core.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *core.User
}

// GuildMemberTypingEvent indicates that a core.Member started typing in a Channel (requires discord.GatewayIntentGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID discord.Snowflake
}

// TextChannel returns the Channel the GuildMemberTypingEvent happened in.
// This will only check cached channels!
func (e GuildMemberTypingEvent) TextChannel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
