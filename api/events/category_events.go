package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericCategoryEvent struct {
	GenericChannelEvent
	Category *api.Category
}

type CategoryCreateEvent struct {
	GenericCategoryEvent
}

type CategoryUpdateEvent struct {
	GenericCategoryEvent
	OldCategory *api.Category
}

type CategoryDeleteEvent struct {
	GenericCategoryEvent
}
