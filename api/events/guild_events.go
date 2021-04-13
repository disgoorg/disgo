package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildEvent generic api.Guild api.GenericEvent
type GenericGuildEvent struct {
	GenericEvent
	Guild *api.Guild
}

// GuildUpdateEvent called upon receiving api.Guild updates
type GuildUpdateEvent struct {
	GenericGuildEvent
	OldGuild *api.Guild
}

// GuildAvailableEvent called when an unavailable api.Guild becomes available
type GuildAvailableEvent struct {
	GenericGuildEvent
}

// GuildUnavailableEvent called when an available api.Guild becomes unavailable
type GuildUnavailableEvent struct {
	GenericGuildEvent
}

// GuildJoinEvent called when the bot joins a api.Guild
type GuildJoinEvent struct {
	GenericGuildEvent
}

// GuildLeaveEvent called when the bot leaves a api.Guild
type GuildLeaveEvent struct {
	GenericGuildEvent
}

// GuildReadyEvent called when the loaded the api.Guild in login phase
type GuildReadyEvent struct {
	GenericGuildEvent
}

type GuildBanEvent struct {
	GenericGuildEvent
	User *api.User
}

type GuildUnbanEvent struct {
	GenericGuildEvent
	User *api.User
}
