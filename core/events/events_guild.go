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

// GuildUpdateEvent is called upon receiving core.Guild updates
type GuildUpdateEvent struct {
	*GenericGuildEvent
	OldGuild *core.Guild
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

// GuildReadyEvent is called when a core.Guild becomes loaded for the first time
type GuildReadyEvent struct {
	*GenericGuildEvent
}

// GuildsReadyEvent is called when all core.Guild(s) are loaded after logging in
type GuildsReadyEvent struct {
	*GenericEvent
	ShardID int
}

// GuildBanEvent is called when a core.Member/core.User is banned from the core.Guild
type GuildBanEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
	User    *core.User
}

// GuildUnbanEvent is called when a core.Member/core.User is unbanned from the core.Guild
type GuildUnbanEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
	User    *core.User
}
