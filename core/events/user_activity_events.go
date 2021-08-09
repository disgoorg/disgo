package events

import (
	
)

// GenericUserActivityEvent is called upon receiving UserActivityStartEvent, UserActivityUpdateEvent or UserActivityEndEvent(requires the api.GatewayIntentsGuildPresences)
type GenericUserActivityEvent struct {
	*GenericGuildMemberEvent
	Member *core.Member
}

// UserActivityStartEvent indicates that an api.User started a new api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityStartEvent struct {
	*GenericUserActivityEvent
	Activity *core.Activity
}

// UserActivityUpdateEvent indicates that an api.User's api.Activity(s) updated(requires the api.GatewayIntentsGuildPresences)
type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	NewActivities *core.Activity
	OldActivities *core.Activity
}

// UserActivityEndEvent indicates that an api.User ended an api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityEndEvent struct {
	*GenericUserActivityEvent
	Activity *core.Activity
}
