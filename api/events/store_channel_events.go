package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericStoreChannelEvent struct {
	GenericChannelEvent
	StoreChannel *api.StoreChannel
}

type StoreChannelCreateEvent struct {
	GenericStoreChannelEvent
}

type StoreChannelUpdateEvent struct {
	GenericStoreChannelEvent
	OldStoreChannel *api.StoreChannel
}

type StoreChannelDeleteEvent struct {
	GenericStoreChannelEvent
}
