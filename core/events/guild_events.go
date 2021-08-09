package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildEvent is called upon receiving GuildUpdateEvent, GuildAvailableEvent, GuildUnavailableEvent, GuildJoinEvent, GuildLeaveEvent, GuildReadyEvent, GuildBanEvent, GuildUnbanEvent
type GenericGuildEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
	Guild   *core.Guild
}

// GuildUpdateEvent is called upon receiving api.Guild updates
type GuildUpdateEvent struct {
	*GenericGuildEvent
	OldGuild *core.Guild
}

// GuildAvailableEvent is called when an unavailable api.Guild becomes available
type GuildAvailableEvent struct {
	*GenericGuildEvent
}

// GuildUnavailableEvent is called when an available api.Guild becomes unavailable
type GuildUnavailableEvent struct {
	*GenericGuildEvent
}

// GuildJoinEvent is called when the bot joins an api.Guild
type GuildJoinEvent struct {
	*GenericGuildEvent
}

// GuildLeaveEvent is called when the bot leaves an api.Guild
type GuildLeaveEvent struct {
	*GenericGuildEvent
}

// GuildReadyEvent is called when the bot loaded the api.Guild in login phase
type GuildReadyEvent struct {
	*GenericGuildEvent
}

// GuildBanEvent is called when an api.Member/api.User is banned from the api.Guild
type GuildBanEvent struct {
	*GenericGuildEvent
	User *core.User
}

// GuildUnbanEvent is called when an api.Member/api.User is unbanned from the api.Guild
type GuildUnbanEvent struct {
	*GenericGuildEvent
	User *core.User
}
