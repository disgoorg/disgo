package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericUserActivityEvent is called upon receiving UserActivityStartEvent, UserActivityUpdateEvent or UserActivityEndEvent(requires the api.GatewayIntentsGuildPresences)
type GenericUserActivityEvent struct {
	*GenericGuildMemberEvent
	Member *Member
}

// UserActivityStartEvent indicates that an api.User started a new api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityStartEvent struct {
	*GenericUserActivityEvent
	Activity discord.Activity
}

// UserActivityUpdateEvent indicates that an api.User's api.Activity(s) updated(requires the api.GatewayIntentsGuildPresences)
type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	NewActivities discord.Activity
	OldActivities discord.Activity
}

// UserActivityEndEvent indicates that an api.User ended an api.Activity(requires the api.GatewayIntentsGuildPresences)
type UserActivityEndEvent struct {
	*GenericUserActivityEvent
	Activity discord.Activity
}
