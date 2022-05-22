package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GenericThread struct {
	*GenericEvent
	Thread   discord.GuildThread
	ThreadID snowflake.ID
	GuildID  snowflake.ID
	ParentID snowflake.ID
}

type ThreadCreate struct {
	*GenericThread
	ThreadMember discord.ThreadMember
}

type ThreadUpdate struct {
	*GenericThread
	OldThread discord.GuildThread
}

type ThreadDelete struct {
	*GenericThread
}

type ThreadShow struct {
	*GenericThread
}

type ThreadHide struct {
	*GenericThread
}

type GenericThreadMember struct {
	*GenericEvent
	GuildID        snowflake.ID
	ThreadID       snowflake.ID
	ThreadMemberID snowflake.ID
	ThreadMember   discord.ThreadMember
}

type ThreadMemberAdd struct {
	*GenericThreadMember
	Member   discord.Member
	Presence *discord.Presence
}

type ThreadMemberUpdate struct {
	*GenericThreadMember
	OldThreadMember discord.ThreadMember
}

type ThreadMemberRemove struct {
	*GenericThreadMember
}
