package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildMemberEvent generic api.Member event
type GenericGuildMemberEvent struct {
	*GenericGuildEvent
	Member *api.Member
}

// User gets the api.User form the api.Cache
func (e GenericGuildMemberEvent) User() *api.User {
	if e.Member == nil {
		return nil
	}
	return e.Disgo().Cache().User(e.Member.User.ID)
}

// GuildMemberJoinEvent indicates that a api.Member joined the api.Guild
type GuildMemberJoinEvent struct {
	*GenericGuildMemberEvent
}

// GuildMemberUpdateEvent indicates that a api.Member updated
type GuildMemberUpdateEvent struct {
	*GenericGuildMemberEvent
	OldMember *api.Member
}

// GuildMemberLeaveEvent indicates that a api.Member left the api.Guild
type GuildMemberLeaveEvent struct {
	*GenericGuildMemberEvent
	User *api.User
}

// GuildMemberTypingEvent indicates that a api.Member started typing in a api.TextChannel(requires api.GatewayIntentsGuildMessageTyping)
type GuildMemberTypingEvent struct {
	*GenericGuildMemberEvent
	ChannelID api.Snowflake
}

// TextChannel returns the api.TextChannel the GuildMemberTypingEvent happened in
func (e GuildMemberTypingEvent) TextChannel() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}
