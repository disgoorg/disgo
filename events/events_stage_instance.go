package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericStageInstanceEvent struct {
	*GenericEvent
	StageInstanceID discord.Snowflake
	StageInstance   *core.StageInstance
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
