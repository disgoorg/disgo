package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type GenericThreadEvent struct {
	*GenericEvent
	Thread   discord.GuildThread
	ThreadID snowflake.Snowflake
	GuildID  snowflake.Snowflake
	ParentID snowflake.Snowflake
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
	GuildID        snowflake.Snowflake
	ThreadID       snowflake.Snowflake
	ThreadMemberID snowflake.Snowflake
	ThreadMember   discord.ThreadMember
}

type ThreadMemberAddEvent struct {
	*GenericThreadMemberEvent
}

type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember discord.ThreadMember
}

type ThreadMemberRemoveEvent struct {
	*GenericThreadMemberEvent
}
