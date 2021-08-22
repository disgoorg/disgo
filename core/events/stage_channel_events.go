package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// GenericStageChannelEvent is called upon receiving StageChannelCreateEvent, StageChannelUpdateEvent or StageChannelDeleteEvent
type GenericStageChannelEvent struct {
	*GenericGuildChannelEvent
	StageChannel core.StageChannel
}

// StageChannelCreateEvent indicates that a new api.StageChannel got created in an api.Guild
type StageChannelCreateEvent struct {
	*GenericStageChannelEvent
}

// StageChannelUpdateEvent indicates that an api.StageChannel got updated in an api.Guild
type StageChannelUpdateEvent struct {
	*GenericStageChannelEvent
	OldStageChannel core.StageChannel
}

// StageChannelDeleteEvent indicates that an api.StageChannel got deleted in an api.Guild
type StageChannelDeleteEvent struct {
	*GenericStageChannelEvent
}
