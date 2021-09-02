package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ GuildService = (*GuildServiceImpl)(nil)

func NewGuildService(restClient Client) GuildService {
	return &GuildServiceImpl{restClient: restClient}
}

type GuildService interface {
	Service
	GetGuild(guildID discord.Snowflake, withCounts bool, opts ...RequestOpt) (*discord.Guild, Error)
	GetGuildPreview(guildID discord.Snowflake, opts ...RequestOpt) (*discord.GuildPreview, Error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.Guild, Error)
	UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.Guild, Error)
	DeleteGuild(guildID discord.Snowflake, opts ...RequestOpt) Error

	GetRoles(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Role, Error)
	CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, Error)
	UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, Error)
	UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, Error)
	DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error

	GetMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Member, Error)
	GetMembers(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Member, Error)
	SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...RequestOpt) ([]discord.Member, Error)
	AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, Error)
	RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error
	UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, Error)
	AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error
	RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error

	UpdateSelfNick(guildID discord.Snowflake, nick string, opts ...RequestOpt) (*string, Error)

	GetBans(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Ban, Error)
	GetBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Ban, Error)
	AddBan(guildID discord.Snowflake, userID discord.Snowflake, deleteMessageDays int, opts ...RequestOpt) Error
	DeleteBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error

	GetIntegrations(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Integration, Error)
	DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake, opts ...RequestOpt) Error
}

type GuildServiceImpl struct {
	restClient Client
}

func (s *GuildServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *GuildServiceImpl) GetGuild(guildID discord.Snowflake, withCounts bool, opts ...RequestOpt) (guild *discord.Guild, rErr Error) {
	values := route.QueryValues{}
	if withCounts {
		values["withCounts"] = true
	}
	compiledRoute, err := route.GetGuild.Compile(values, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guild, opts...)
	return
}

func (s *GuildServiceImpl) GetGuildPreview(guildID discord.Snowflake, opts ...RequestOpt) (guildPreview *discord.GuildPreview, rErr Error) {
	compiledRoute, err := route.GetGuildPreview.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guildPreview, opts...)
	return
}

func (s *GuildServiceImpl) CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (guild *discord.Guild, rErr Error) {
	compiledRoute, err := route.CreateGuild.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, guildCreate, &guild, opts...)
	return
}

func (s *GuildServiceImpl) UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (guild *discord.Guild, rErr Error) {
	compiledRoute, err := route.UpdateGuild.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, guildUpdate, &guild, opts...)
	return
}

func (s *GuildServiceImpl) DeleteGuild(guildID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) GetRoles(guildID discord.Snowflake, opts ...RequestOpt) (roles []discord.Role, rErr Error) {
	compiledRoute, err := route.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &roles, opts...)
	return
}

func (s *GuildServiceImpl) CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (role *discord.Role, rErr Error) {
	compiledRoute, err := route.CreateRole.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, createRole, &role, opts...)
	return
}

func (s *GuildServiceImpl) UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (role *discord.Role, rErr Error) {
	compiledRoute, err := route.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, roleUpdate, &role, opts...)
	return
}

func (s *GuildServiceImpl) UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) (roles []discord.Role, rErr Error) {
	compiledRoute, err := route.UpdateRolePositions.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, rolePositionUpdates, &roles, opts...)
	return
}

func (s *GuildServiceImpl) DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteRole.Compile(nil, guildID, roleID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) GetMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (member *discord.Member, rErr Error) {
	compiledRoute, err := route.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &member, opts...)
	return
}

func (s *GuildServiceImpl) GetMembers(guildID discord.Snowflake, opts ...RequestOpt) (members []discord.Member, rErr Error) {
	compiledRoute, err := route.GetMembers.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *GuildServiceImpl) SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...RequestOpt) (members []discord.Member, rErr Error) {
	values := route.QueryValues{}
	if query != "" {
		values["query"] = query
	}
	if limit != 0 {
		values["limit"] = limit
	}
	compiledRoute, err := route.SearchMembers.Compile(values, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *GuildServiceImpl) AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (member *discord.Member, rErr Error) {
	compiledRoute, err := route.GetMembers.Compile(nil, guildID, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, memberAdd, &member, opts...)
	return
}

func (s *GuildServiceImpl) RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveMember.Compile(nil, guildID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (member *discord.Member, rErr Error) {
	compiledRoute, err := route.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, memberUpdate, &member, opts...)
	return
}

func (s *GuildServiceImpl) AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) UpdateSelfNick(guildID discord.Snowflake, nick string, opts ...RequestOpt) (nickName *string, rErr Error) {
	compiledRoute, err := route.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, discord.SelfNickUpdate{Nick: nick}, nickName, opts...)
	return
}

func (s *GuildServiceImpl) GetBans(guildID discord.Snowflake, opts ...RequestOpt) (bans []discord.Ban, rErr Error) {
	compiledRoute, err := route.GetBans.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &bans, opts...)
	return
}

func (s *GuildServiceImpl) GetBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (ban *discord.Ban, rErr Error) {
	compiledRoute, err := route.GetBan.Compile(nil, guildID, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &ban, opts...)
	return
}

func (s *GuildServiceImpl) AddBan(guildID discord.Snowflake, userID discord.Snowflake, deleteMessageDays int, opts ...RequestOpt) Error {
	compiledRoute, err := route.AddBan.Compile(nil, guildID, userID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, discord.AddBan{DeleteMessageDays: deleteMessageDays}, nil, opts...)
}

func (s *GuildServiceImpl) DeleteBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteBan.Compile(nil, guildID, userID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) GetIntegrations(guildID discord.Snowflake, opts ...RequestOpt) (integrations []discord.Integration, rErr Error) {
	compiledRoute, err := route.GetIntegrations.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &integrations, opts...)
	return
}

func (s *GuildServiceImpl) DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteIntegration.Compile(nil, guildID, integrationID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *GuildServiceImpl) GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, rErr Error) {
	compiledRoute, err := route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}
