package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

type GenericThreadEvent struct {
	*GenericEvent
	Thread   core.GuildThread
	ThreadID snowflake.Snowflake
	GuildID  snowflake.Snowflake
	ParentID snowflake.Snowflake
}

type ThreadCreateEvent struct {
	*GenericThreadEvent
}

type ThreadUpdateEvent struct {
	*GenericThreadEvent
	OldThread core.GuildThread
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
	ThreadMember   *core.ThreadMember
}

type ThreadMemberAddEvent struct {
	*GenericThreadMemberEvent
}

type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember *core.ThreadMember
}

type ThreadMemberRemoveEvent struct {
	*GenericThreadMemberEvent
}
