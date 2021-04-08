package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericUserActivityEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}

type UserActivityStartEvent struct {
	GenericUserActivityEvent
	Activity *api.Activity
}

type UserActivityUpdateEvent struct {
	GenericUserActivityEvent
	NewActivities *api.Activity
	OldActivities *api.Activity
}

type UserActivityEndEvent struct {
	GenericUserActivityEvent
	Activity *api.Activity
}
