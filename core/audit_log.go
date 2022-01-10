package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type AuditLog struct {
	discord.AuditLog
	GuildScheduledEvents map[discord.Snowflake]*GuildScheduledEvent
	Integrations         map[discord.Snowflake]Integration
	Threads              map[discord.Snowflake]GuildThread
	Users                map[discord.Snowflake]*User
	Webhooks             map[discord.Snowflake]Webhook
	GuildID              discord.Snowflake
	FilterOptions        AuditLogFilterOptions
	Bot                  *Bot
}

func (l *AuditLog) Guild() (Guild, bool) {
	return l.Bot.Caches.Guilds().Get(l.GuildID)
}

// AuditLogFilterOptions fields used to filter audit-log retrieving
type AuditLogFilterOptions struct {
	UserID     discord.Snowflake
	ActionType discord.AuditLogEvent
	Before     discord.Snowflake
	Limit      int
}

// Before gets new AuditLog(s) from Discord before the last one
func (l *AuditLog) Before(opts ...rest.RequestOpt) (*AuditLog, error) {
	before := discord.Snowflake("")
	if len(l.Entries) > 0 {
		before = l.Entries[len(l.Entries)-1].ID
	}
	auditLog, err := l.Bot.RestServices.AuditLogService().GetAuditLog(l.GuildID, l.FilterOptions.UserID, l.FilterOptions.ActionType, before, l.FilterOptions.Limit, opts...)
	if err != nil {
		return nil, err
	}
	return l.Bot.EntityBuilder.CreateAuditLog(l.GuildID, *auditLog, l.FilterOptions, CacheStrategyNoWs), nil
}
