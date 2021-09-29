package core

// GenericStageInstanceEvent generic StageInstance event
type GenericStageInstanceEvent struct {
	*GenericGuildChannelEvent
	StageInstance *StageInstance
}

// GenericStageInstanceEvent indicates that a StageInstance got created
type StageInstanceCreateEvent struct {
	*GenericStageInstanceEvent
}

// StageInstanceUpdateEvent indicates that a StageInstance got updated
type StageInstanceUpdateEvent struct {
	*GenericStageInstanceEvent
	OldStageInstance *StageInstance
}

// StageInstanceDeleteEvent indicates that a StageInstance got deleted
type StageInstanceDeleteEvent struct {
	*GenericStageInstanceEvent
}
