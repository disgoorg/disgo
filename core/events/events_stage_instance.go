package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericStageInstanceEvent generic StageInstance event
type GenericStageInstanceEvent struct {
	*GenericEvent
	StageInstanceID snowflake.Snowflake
	StageInstance   discord.StageInstance
}

// StageInstanceCreateEvent indicates that a StageInstance got created
type StageInstanceCreateEvent struct {
	*GenericStageInstanceEvent
}

// StageInstanceUpdateEvent indicates that a StageInstance got updated
type StageInstanceUpdateEvent struct {
	*GenericStageInstanceEvent
	OldStageInstance discord.StageInstance
}

// StageInstanceDeleteEvent indicates that a StageInstance got deleted
type StageInstanceDeleteEvent struct {
	*GenericStageInstanceEvent
}
