package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuild is called upon receiving GuildUpdate , GuildAvailable , GuildUnavailable , GuildJoin , GuildLeave , GuildReady , GuildBan , GuildUnban
type GenericGuild struct {
	*GenericEvent
	GuildID snowflake.ID
	Guild   discord.Guild
}

// GuildUpdate is called upon receiving discord.Guild updates
type GuildUpdate struct {
	*GenericGuild
	OldGuild discord.Guild
}

// GuildAvailable is called when an unavailable discord.Guild becomes available
type GuildAvailable struct {
	*GenericGuild
}

// GuildUnavailable is called when an available discord.Guild becomes unavailable
type GuildUnavailable struct {
	*GenericGuild
}

// GuildJoin is called when the bot joins a discord.Guild
type GuildJoin struct {
	*GenericGuild
}

// GuildLeave is called when the bot leaves a discord.Guild
type GuildLeave struct {
	*GenericGuild
}

// GuildReady is called when a discord.Guild becomes loaded for the first time
type GuildReady struct {
	*GenericGuild
}

// GuildsReady is called when all discord.Guild(s) are loaded after logging in
type GuildsReady struct {
	*GenericEvent
}

// GuildBan is called when a discord.Member/discord.User is banned from the discord.Guild
type GuildBan struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
}

// GuildUnban is called when a discord.Member/discord.User is unbanned from the discord.Guild
type GuildUnban struct {
	*GenericEvent
	GuildID snowflake.ID
	User    discord.User
}

// GuildAuditLogEntryCreate is called when a new discord.AuditLogEntry is created
type GuildAuditLogEntryCreate struct {
	*GenericEvent
	GuildID       snowflake.ID
	AuditLogEntry discord.AuditLogEntry
}
