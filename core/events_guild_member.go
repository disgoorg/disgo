package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMemberEvent generic api.Member event
type GenericGuildMemberEvent struct {
	*GenericGuildEvent
	Member *Member
}

// GuildMemberJoinEvent indicates that an api.Member joined the api.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that an api.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *Member
}

// GuildMemberLeaveEvent indicates that an api.Member left the api.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *User
}

// GuildMemberTypingEvent indicates that an api.Member started typing in an api.TextChannel(requires api.GatewayIntentsGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID discord.Snowflake
}

// TextChannel returns the api.TextChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) TextChannel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}
