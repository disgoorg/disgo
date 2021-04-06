package events

import "github.com/DisgoOrg/disgo/api"

// GenericGuildMemberEvent generic api.Member event
type GenericGuildMemberEvent struct {
	GenericGuildEvent
	UserID api.Snowflake
}

// User gets the api.User form the api.Cache
func (e GenericGuildMemberEvent) User() *api.User {
	return e.Disgo().Cache().User(e.UserID)
}

// GuildMemberJoinEvent indicates that a api.Member joined the api.Guild
type GuildMemberJoinEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}

// GuildMemberUpdateEvent indicates that a api.Member updated
type GuildMemberUpdateEvent struct {
	GenericGuildMemberEvent
	OldMember *api.Member
	NewMember *api.Member
}

// GuildMemberLeaveEvent indicates that a api.Member left the api.Guild
type GuildMemberLeaveEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}
