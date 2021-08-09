package events

import (
	
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildMemberEvent generic api.Member event
type GenericGuildMemberEvent struct {
	*GenericGuildEvent
	Member *core.Member
}

// User gets the api.User form the api.Cache
func (e GenericGuildMemberEvent) User() *core.User {
	if e.Member == nil {
		return nil
	}
	return e.Disgo().Cache().User(e.Member.User.ID)
}

// GuildMemberJoinEvent indicates that an api.Member joined the api.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that an api.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *core.Member
}

// GuildMemberLeaveEvent indicates that an api.Member left the api.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *core.User
}

// GuildMemberTypingEvent indicates that an api.Member started typing in an api.TextChannel(requires api.GatewayIntentsGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID discord.Snowflake
}

// TextChannel returns the api.TextChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) TextChannel() *core.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}
