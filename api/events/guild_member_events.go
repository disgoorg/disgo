package events

import "github.com/DiscoOrg/disgo/api"

type GenericGuildMemberEvent struct {
	GenericGuildEvent
	UserID api.Snowflake
}

func (e GenericGuildMemberEvent) User() *api.User {
	return e.Disgo.Cache().User(e.UserID)
}

type GuildMemberJoinEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}

type GuildMemberUpdateEvent struct {
	GenericGuildMemberEvent
	OldMember *api.Member
	NewMember *api.Member
}

type GuildMemberLeaveEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}
