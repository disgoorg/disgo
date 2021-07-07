package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericUserActivityEvent is called upon receiving UserActivityStartEvent, UserActivityUpdateEvent or UserActivityEndEvent(requires the api.GatewayIntentsGuildPresences)
type GenericUserActivityEvent struct {
	*GenericGuildMemberEvent
	Member *api.Member
}

// UserActivityStartEvent indicates that a api.User started a new api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityStartEvent struct {
	*GenericUserActivityEvent
	Activity *api.Activity
}

// UserActivityUpdateEvent indicates that a api.User's api.Activity(s) updated(requires the api.GatewayIntentsGuildPresences)
type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	NewActivities *api.Activity
	OldActivities *api.Activity
}

// UserActivityEndEvent indicates that a api.User ended a api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityEndEvent struct {
	*GenericUserActivityEvent
	Activity *api.Activity
}
