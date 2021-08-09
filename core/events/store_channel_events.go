package events

import (
	
)

// GenericStoreChannelEvent is called upon receiving StoreChannelCreateEvent, StoreChannelUpdateEvent or StoreChannelDeleteEvent
type GenericStoreChannelEvent struct {
	*GenericGuildChannelEvent
	StoreChannel *core.StoreChannel
}

// StoreChannelCreateEvent indicates that a new api.StoreChannel got created in an api.Guild
type StoreChannelCreateEvent struct {
	*GenericStoreChannelEvent
}

// StoreChannelUpdateEvent indicates that an api.StoreChannel got updated in an api.Guild
type StoreChannelUpdateEvent struct {
	*GenericStoreChannelEvent
	OldStoreChannel *core.StoreChannel
}

// StoreChannelDeleteEvent indicates that an api.StoreChannel got deleted in an api.Guild
type StoreChannelDeleteEvent struct {
	*GenericStoreChannelEvent
}
