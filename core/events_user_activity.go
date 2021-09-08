package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericUserActivityEvent is called upon receiving UserActivityStartEvent, UserActivityUpdateEvent or UserActivityEndEvent(requires the core.GatewayIntentsGuildPresences)
type GenericUserActivityEvent struct {
	*GenericGuildMemberEvent
	Member *Member
}

// UserActivityStartEvent indicates that an core.User started a new core.Activity(requires the core.GatewayIntentsGuildPresences)
type UserActivityStartEvent struct {
	*GenericUserActivityEvent
	Activity discord.Activity
}

// UserActivityUpdateEvent indicates that an core.User's core.Activity(s) updated(requires the core.GatewayIntentsGuildPresences)
type UserActivityUpdateEvent struct {
	*GenericUserActivityEvent
	NewActivities discord.Activity
	OldActivities discord.Activity
}

// UserActivityEndEvent indicates that an core.User ended an core.Activity(requires the core.GatewayIntentsGuildPresences)
type UserActivityEndEvent struct {
	*GenericUserActivityEvent
	Activity discord.Activity
}
