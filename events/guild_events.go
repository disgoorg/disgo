package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildEvent is called upon receiving GuildUpdateEvent, GuildAvailableEvent, GuildUnavailableEvent, GuildJoinEvent, GuildLeaveEvent, GuildReadyEvent, GuildBanEvent, GuildUnbanEvent
type GenericGuildEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	Guild   discord.Guild
}

// GuildUpdateEvent is called upon receiving discord.Guild updates
type GuildUpdateEvent struct {
	*GenericGuildEvent
	OldGuild discord.Guild
}

// GuildAvailableEvent is called when an unavailable discord.Guild becomes available
type GuildAvailableEvent struct {
	*GenericGuildEvent
}

// GuildUnavailableEvent is called when an available discord.Guild becomes unavailable
type GuildUnavailableEvent struct {
	*GenericGuildEvent
}

// GuildJoinEvent is called when the bot joins a discord.Guild
type GuildJoinEvent struct {
	*GenericGuildEvent
}

// GuildLeaveEvent is called when the bot leaves a discord.Guild
type GuildLeaveEvent struct {
	*GenericGuildEvent
}

// GuildReadyEvent is called when a discord.Guild becomes loaded for the first time
type GuildReadyEvent struct {
	*GenericGuildEvent
}

// GuildsReadyEvent is called when all discord.Guild(s) are loaded after logging in
type GuildsReadyEvent struct {
	*GenericEvent
	ShardID int
}

// GuildBanEvent is called when a discord.Member/discord.User is banned from the discord.Guild
type GuildBanEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
}

// GuildUnbanEvent is called when a discord.Member/discord.User is unbanned from the discord.Guild
type GuildUnbanEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
}
