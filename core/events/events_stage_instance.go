package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericStageInstanceEvent generic StageInstance event
type GenericStageInstanceEvent struct {
	*GenericEvent
	StageInstanceID discord.Snowflake
	StageInstance   *core.StageInstance
}

// StageInstanceCreateEvent indicates that a StageInstance got created
type StageInstanceCreateEvent struct {
	*GenericStageInstanceEvent
}

// StageInstanceUpdateEvent indicates that a StageInstance got updated
type StageInstanceUpdateEvent struct {
	*GenericStageInstanceEvent
	OldStageInstance *core.StageInstance
}

// StageInstanceDeleteEvent indicates that a StageInstance got deleted
type StageInstanceDeleteEvent struct {
	*GenericStageInstanceEvent
}
