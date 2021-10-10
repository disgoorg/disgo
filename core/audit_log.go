package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type AuditLog struct {
	discord.AuditLog
	Bot           *Bot
	GuildID       discord.Snowflake
	Users         map[discord.Snowflake]*User
	Integrations  map[discord.Snowflake]*Integration
	Webhooks      map[discord.Snowflake]*Webhook
	FilterOptions AuditLogFilterOptions
}

func (l *AuditLog) Guild() *Guild {
	return l.Bot.Caches.GuildCache().Get(l.GuildID)
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
