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

// GuildJoinEvent is called when the bot joins an core.Guild
type GuildJoinEvent struct {
	*GenericGuildEvent
}

// GuildLeaveEvent is called when the bot leaves an core.Guild
type GuildLeaveEvent struct {
	*GenericGuildEvent
}

// GuildBanEvent is called when an core.Member/core.User is banned from the core.Guild
type GuildBanEvent struct {
	*GenericGuildEvent
	User *core.User
}

// GuildUnbanEvent is called when an core.Member/core.User is unbanned from the core.Guild
type GuildUnbanEvent struct {
	*GenericGuildEvent
	User *core.User
}
