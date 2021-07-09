package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericThreadEvent is called upon receiving ThreadCreateEvent, ThreadUpdateEvent or ThreadDeleteEvent
type GenericThreadEvent struct {
	*GenericChannelEvent
	Thread api.Thread
}

// ThreadCreateEvent indicates that a new api.Thread got created in a api.Guild
type ThreadCreateEvent struct {
	*GenericThreadEvent
}

// ThreadUpdateEvent indicates that a api.Thread got updated in a api.Guild
type ThreadUpdateEvent struct {
	*GenericThreadEvent
	OldThread api.Thread
}

// ThreadDeleteEvent indicates that a api.Thread got deleted in a api.Guild
type ThreadDeleteEvent struct {
	*GenericThreadEvent
}

// ThreadJoinEvent indicates you joined a api.Thread
type ThreadJoinEvent struct {
	*GenericThreadEvent
}

// ThreadLeaveEvent indicates you left a api.Thread
type ThreadLeaveEvent struct {
	*GenericThreadEvent
}
