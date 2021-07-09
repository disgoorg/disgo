package events

import "github.com/DisgoOrg/disgo/api"

// GenericThreadMemberEvent is called upon receiving ThreadMemberAddEvent, ThreadMemberUpdateEvent or ThreadMemberRemoveEvent
type GenericThreadMemberEvent struct {
	*GenericThreadEvent
	ThreadMember *api.ThreadMember
}

// ThreadMemberAddEvent indicates that a api.ThreadMember joined the api.Thread
type ThreadMemberAddEvent struct {
	*GenericThreadMemberEvent
}

// ThreadMemberUpdateEvent indicates that a api.ThreadMember updated
type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember *api.ThreadMember
}

// ThreadMemberRemoveEvent indicates that a api.ThreadMember left the api.Thread
type ThreadMemberRemoveEvent struct {
	*GenericThreadMemberEvent
}
