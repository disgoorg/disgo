package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// GenericNewsChannelEvent is called upon receiving NewsChannelCreateEvent, NewsChannelUpdateEvent or NewsChannelDeleteEvent
type GenericNewsChannelEvent struct {
	*GenericGuildChannelEvent
	NewsChannel core.NewsChannel
}

// NewsChannelCreateEvent indicates that a new core.NewsChannel got created in a core.Guild
type NewsChannelCreateEvent struct {
	*GenericNewsChannelEvent
}

// NewsChannelUpdateEvent indicates that a core.NewsChannel got updated in a core.Guild
type NewsChannelUpdateEvent struct {
	*GenericNewsChannelEvent
	OldNewsChannel core.NewsChannel
}

// NewsChannelDeleteEvent indicates that a core.NewsChannel got deleted in a core.Guild
type NewsChannelDeleteEvent struct {
	*GenericNewsChannelEvent
}
