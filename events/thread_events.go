package events

import "github.com/DisgoOrg/disgo/core"

type GenericThreadEvent struct {
	*GenericGuildChannelEvent
	Thread core.GuildThread
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

type GenericThreadMemberEvent struct {
	*GenericThreadEvent
	ThreadMember *core.ThreadMember
}

type ThreadMemberJoinEvent struct {
	*GenericThreadMemberEvent
}

type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember *core.ThreadMember
}

type ThreadMemberLeaveEvent struct {
	*GenericThreadMemberEvent
}
