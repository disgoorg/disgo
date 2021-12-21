package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericThreadEvent struct {
	*GenericEvent
	Thread   core.GuildThread
	ThreadID discord.Snowflake
	GuildID  discord.Snowflake
	ParentID discord.Snowflake
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
	GuildID        discord.Snowflake
	ThreadID       discord.Snowflake
	ThreadMemberID discord.Snowflake
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
