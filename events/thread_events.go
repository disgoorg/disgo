package events

type ThreadCreateEvent struct {
	*GenericGuildChannelEvent
}

type ThreadUpdateEvent struct {
	*GenericGuildChannelEvent
}

type ThreadDeleteEvent struct {
	*GenericGuildChannelEvent
}

type ThreadJoinEvent struct {
}

type ThreadLeaveEvent struct {
}
