package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewAuditLogService(client Client) AuditLogService {
	return nil
}

type AuditLogService interface {
	Service
	GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int) (*discord.AuditLog, Error)
}
