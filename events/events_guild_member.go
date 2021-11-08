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

// GuildMemberJoinEvent indicates that an core.Member joined the core.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that an core.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *core.Member
}

// GuildMemberLeaveEvent indicates that an core.Member left the core.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *core.User
}

// GuildMemberTypingEvent indicates that an core.Member started typing in an core.TextChannel(requires core.GatewayIntentsGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID discord.Snowflake
}

// MessageChannel returns the core.GuildTextChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) MessageChannel() core.GuildMessageChannel {
	if ch := e.Bot().Caches.ChannelCache().Get(e.ChannelID); ch != nil {
		return ch.(core.GuildMessageChannel)
	}
	return nil
}
