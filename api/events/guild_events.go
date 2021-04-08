package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildEvent generic api.Guild api.GenericEvent
type GenericGuildEvent struct {
	GenericEvent
	GuildID api.Snowflake
}

// Guild returns the api.Guild from the api.Cache
func (e GenericGuildEvent) Guild() *api.Guild {
	return e.Disgo().Cache().Guild(e.GuildID)
}

// GuildUpdateEvent called upon receiving api.Guild updates
type GuildUpdateEvent struct {
	GenericGuildEvent
	NewGuild *api.Guild
	OldGuild *api.Guild
}

// GuildAvailableEvent called when an unavailable api.Guild becomes available
type GuildAvailableEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

// GuildUnavailableEvent called when an available api.Guild becomes unavailable
type GuildUnavailableEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

// GuildJoinEvent called when the bot joins a api.Guild
type GuildJoinEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

// GuildLeaveEvent called when the bot leaves a api.Guild
type GuildLeaveEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

// GuildReadyEvent called when the loaded the api.Guild in login phase
type GuildReadyEvent struct {
	GenericGuildEvent
	Guild *api.Guild
}

type GuildBanEvent struct {
	GenericGuildEvent
	User *api.User
}

type GuildUnbanEvent struct {
	GenericGuildEvent
	User *api.User
}
