package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericStageInstance generic StageInstance event
type GenericStageInstance struct {
	*GenericEvent
	StageInstanceID snowflake.ID
	StageInstance   discord.StageInstance
}

// StageInstanceCreate indicates that a StageInstance got created
type StageInstanceCreate struct {
	*GenericStageInstance
}

// StageInstanceUpdate indicates that a StageInstance got updated
type StageInstanceUpdate struct {
	*GenericStageInstance
	OldStageInstance discord.StageInstance
}

// StageInstanceDelete indicates that a StageInstance got deleted
type StageInstanceDelete struct {
	*GenericStageInstance
}
