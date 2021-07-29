package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericCategoryEvent is called upon receiving CategoryCreateEvent, CategoryUpdateEvent or CategoryDeleteEvent
type GenericCategoryEvent struct {
	*GenericGuildChannelEvent
	Category *api.Category
}

// CategoryCreateEvent indicates that a new api.Category got created in an api.Guild
type CategoryCreateEvent struct {
	*GenericCategoryEvent
}

// CategoryUpdateEvent indicates that an api.Category got updated in an api.Guild
type CategoryUpdateEvent struct {
	*GenericCategoryEvent
	OldCategory *api.Category
}

// CategoryDeleteEvent indicates that an api.Category got deleted in an api.Guild
type CategoryDeleteEvent struct {
	*GenericCategoryEvent
}
