package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericCategoryEvent struct {
	GenericChannelEvent
}

func (e GenericCategoryEvent) Category() *api.Category {
	return e.Disgo().Cache().Category(e.ChannelID)
}

type CategoryCreateEvent struct {
	GenericCategoryEvent
	Category *api.Category
}

type CategoryUpdateEvent struct {
	GenericChannelEvent
	NewCategory *api.Category
	OldCategory *api.Category
}

type CategoryDeleteEvent struct {
	GenericChannelEvent
	Category *api.Category
}
