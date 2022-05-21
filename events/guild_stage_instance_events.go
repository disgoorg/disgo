package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericStageInstanceEvent generic StageInstance event
type GenericStageInstanceEvent struct {
	*GenericEvent
	StageInstanceID snowflake.ID
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
