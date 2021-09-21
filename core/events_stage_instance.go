package core

type GenericStageInstanceEvent struct {
	*GenericGuildChannelEvent
	StageInstance *StageInstance
}

type StageInstanceCreateEvent struct {
	*GenericStageInstanceEvent
}

type StageInstanceUpdateEvent struct {
	*GenericStageInstanceEvent
	OldStageInstance *StageInstance
}

type StageInstanceDeleteEvent struct {
	*GenericStageInstanceEvent
}
