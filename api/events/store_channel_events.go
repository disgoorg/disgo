package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericStoreChannelEvent is called upon receiving StoreChannelCreateEvent, StoreChannelUpdateEvent or StoreChannelDeleteEvent
type GenericStoreChannelEvent struct {
	*GenericGuildChannelEvent
	StoreChannel *api.StoreChannel
}

// StoreChannelCreateEvent indicates that a new api.StoreChannel got created in a api.Guild
type StoreChannelCreateEvent struct {
	*GenericStoreChannelEvent
}

// StoreChannelUpdateEvent indicates that a api.StoreChannel got updated in a api.Guild
type StoreChannelUpdateEvent struct {
	*GenericStoreChannelEvent
	OldStoreChannel *api.StoreChannel
}

// StoreChannelDeleteEvent indicates that a api.StoreChannel got deleted in a api.Guild
type StoreChannelDeleteEvent struct {
	*GenericStoreChannelEvent
}
