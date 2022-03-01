package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Service         = (*AuditLogServiceImpl)(nil)
	_ AuditLogService = (*AuditLogServiceImpl)(nil)
)

func NewAuditLogService(restClient Client) AuditLogService {
	return &AuditLogServiceImpl{restClient: restClient}
}

type AuditLogService interface {
	Service
	GetAuditLog(guildID snowflake.Snowflake, userID snowflake.Snowflake, actionType discord.AuditLogEvent, before snowflake.Snowflake, limit int, opts ...RequestOpt) (*discord.AuditLog, error)
}

type AuditLogServiceImpl struct {
	restClient Client
}

func (s *AuditLogServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *AuditLogServiceImpl) GetAuditLog(guildID snowflake.Snowflake, userID snowflake.Snowflake, actionType discord.AuditLogEvent, before snowflake.Snowflake, limit int, opts ...RequestOpt) (auditLog *discord.AuditLog, err error) {
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
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuditLogs.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &auditLog, opts...)
	return
}
