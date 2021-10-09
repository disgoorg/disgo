package events

import "github.com/DisgoOrg/disgo/core"

type GenericStageInstanceEvent struct {
	*GenericGuildChannelEvent
	StageInstance *core.StageInstance
}

type StageInstanceCreateEvent struct {
	*GenericStageInstanceEvent
}

type StageInstanceUpdateEvent struct {
	*GenericStageInstanceEvent
	OldStageInstance *core.StageInstance
}

type StageInstanceDeleteEvent struct {
	*GenericStageInstanceEvent
}
