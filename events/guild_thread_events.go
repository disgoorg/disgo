package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GenericThreadEvent struct {
	*GenericEvent
	Thread   discord.GuildThread
	ThreadID snowflake.ID
	GuildID  snowflake.ID
	ParentID snowflake.ID
}

type ThreadCreateEvent struct {
	*GenericThreadEvent
	ThreadMember discord.ThreadMember
}

type ThreadUpdateEvent struct {
	*GenericThreadEvent
	OldThread discord.GuildThread
}

type ThreadDeleteEvent struct {
	*GenericThreadEvent
}

type ThreadShowEvent struct {
	*GenericThreadEvent
}

type ThreadHideEvent struct {
	*GenericThreadEvent
}

type GenericThreadMemberEvent struct {
	*GenericEvent
	GuildID        snowflake.ID
	ThreadID       snowflake.ID
	ThreadMemberID snowflake.ID
	ThreadMember   discord.ThreadMember
}

type ThreadMemberAddEvent struct {
	*GenericThreadMemberEvent
	Member   discord.Member
	Presence *discord.Presence
}

type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember discord.ThreadMember
}

type ThreadMemberRemoveEvent struct {
	*GenericThreadMemberEvent
}
