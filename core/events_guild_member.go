package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMemberEvent generic core.Member event
type GenericGuildMemberEvent struct {
	*GenericGuildEvent
	Member *Member
}

// GuildMemberJoinEvent indicates that an core.Member joined the core.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that an core.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *Member
}

// GuildMemberLeaveEvent indicates that an core.Member left the core.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *User
}

// GuildMemberTypingEvent indicates that an core.Member started typing in an core.TextChannel(requires core.GatewayIntentsGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID discord.Snowflake
}

// TextChannel returns the core.TextChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) TextChannel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
