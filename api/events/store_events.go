package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericStoreChannelEvent struct {
	GenericChannelEvent
}

func (e GenericStoreChannelEvent) Category() *api.StoreChannel {
	return e.Disgo().Cache().StoreChannel(e.ChannelID)
}

type StoreChannelCreateEvent struct {
	GenericStoreChannelEvent
	StoreChannel *api.StoreChannel
}

type StoreChannelUpdateEvent struct {
	GenericStoreChannelEvent
	NewStoreChannel *api.StoreChannel
	OldStoreChannel *api.StoreChannel
}

type StoreChannelDeleteEvent struct {
	GenericStoreChannelEvent
	StoreChannel *api.StoreChannel
}
