package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildEvent is called upon receiving GuildUpdateEvent, GuildAvailableEvent, GuildUnavailableEvent, GuildJoinEvent, GuildLeaveEvent, GuildReadyEvent, GuildBanEvent, GuildUnbanEvent
type GenericGuildEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
	Guild   *Guild
}

// GuildUpdateEvent is called upon receiving core.Guild updates
type GuildUpdateEvent struct {
	*GenericGuildEvent
	OldGuild *Guild
}

// GuildAvailableEvent is called when an unavailable core.Guild becomes available
type GuildAvailableEvent struct {
	*GenericGuildEvent
}

// GuildUnavailableEvent is called when an available core.Guild becomes unavailable
type GuildUnavailableEvent struct {
	*GenericGuildEvent
}

// GuildJoinEvent is called when the bot joins a core.Guild
type GuildJoinEvent struct {
	*GenericGuildEvent
}

// GuildLeaveEvent is called when the bot leaves a core.Guild
type GuildLeaveEvent struct {
	*GenericGuildEvent
}

// GuildBanEvent is called when a core.Member/core.User is banned from a core.Guild
type GuildBanEvent struct {
	*GenericGuildEvent
	User *User
}

// GuildUnbanEvent is called when a core.Member/core.User is unbanned from a core.Guild
type GuildUnbanEvent struct {
	*GenericGuildEvent
	User *User
}
