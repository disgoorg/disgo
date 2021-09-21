package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ auditLogService = (*AuditLogServiceImpl)(nil)

func NewAuditLogService(restClient Client) auditLogService {
	return &AuditLogServiceImpl{restClient: restClient}
}

type auditLogService interface {
	Service
	GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int, opts ...RequestOpt) (*discord.AuditLog, Error)
}

type AuditLogServiceImpl struct {
	restClient Client
}

func (s *AuditLogServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *AuditLogServiceImpl) GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int, opts ...RequestOpt) (auditLog *discord.AuditLog, rErr Error) {
	values := route.QueryValues{}
	if userID != "" {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != "" {
		values["before"] = guildID
	}
	if limit != 0 {
		values["limit"] = limit
	}
	compiledRoute, err := route.GetAuditLogs.Compile(values, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &auditLog, opts...)
	return
}
