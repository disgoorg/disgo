package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericCategoryEvent is called upon receiving CategoryCreateEvent, CategoryUpdateEvent or CategoryDeleteEvent
type GenericCategoryEvent struct {
	*GenericGuildChannelEvent
	Category *api.Category
}

// CategoryCreateEvent indicates that a new api.Category got created in a api.Guild
type CategoryCreateEvent struct {
	*GenericCategoryEvent
}

// CategoryUpdateEvent indicates that a api.Category got updated in a api.Guild
type CategoryUpdateEvent struct {
	*GenericCategoryEvent
	OldCategory *api.Category
}

// CategoryDeleteEvent indicates that a api.Category got deleted in a api.Guild
type CategoryDeleteEvent struct {
	*GenericCategoryEvent
}
