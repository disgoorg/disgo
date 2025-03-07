package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericGuild represents a generic guild event
type GenericGuild struct {
	*GenericEvent
	GuildID snowflake.ID
}

// GuildUpdate is called upon receiving discord.Guild updates
type GuildUpdate struct {
	*GenericGuild
	Guild    discord.Guild
	OldGuild discord.Guild // the old cached guild
}

// GuildAvailable is called when an unavailable discord.Guild becomes available
type GuildAvailable struct {
	*GenericGuild
	Guild discord.GatewayGuild
}

// GuildUnavailable is called when an available discord.Guild becomes unavailable
type GuildUnavailable struct {
	*GenericGuild
	Guild discord.Guild // the old cached guild
}

// GuildJoin is called when the bot joins a discord.Guild
type GuildJoin struct {
	*GenericGuild
	Guild discord.GatewayGuild
}

// GuildLeave is called when the bot leaves a discord.Guild
type GuildLeave struct {
	*GenericGuild
	Guild discord.Guild // the old cached guild
}

// GuildReady is called when a discord.Guild becomes loaded for the first time
type GuildReady struct {
	*GenericGuild
	Guild discord.GatewayGuild
}

// GuildsReady is called when all discord.Guild(s) are loaded after logging in
type GuildsReady struct {
	*GenericEvent
}

// GuildBan is called when a discord.Member/discord.User is banned from the discord.Guild
type GuildBan struct {
	*GenericGuild
	User discord.User
}

// GuildUnban is called when a discord.Member/discord.User is unbanned from the discord.Guild
type GuildUnban struct {
	*GenericGuild
	User discord.User
}

// GuildAuditLogEntryCreate is called when a new discord.AuditLogEntry is created
type GuildAuditLogEntryCreate struct {
	*GenericGuild
	AuditLogEntry discord.AuditLogEntry
}
