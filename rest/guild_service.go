package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service      = (*guildServiceImpl)(nil)
	_ GuildService = (*guildServiceImpl)(nil)
)

func NewGuildService(restClient Client) GuildService {
	return &guildServiceImpl{restClient: restClient}
}

type GuildService interface {
	Service
	GetGuild(guildID discord.Snowflake, withCounts bool, opts ...RequestOpt) (*discord.Guild, error)
	GetGuildPreview(guildID discord.Snowflake, opts ...RequestOpt) (*discord.GuildPreview, error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.Guild, error)
	UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.Guild, error)
	DeleteGuild(guildID discord.Snowflake, opts ...RequestOpt) error

	GetRoles(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Role, error)
	GetRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) ([]discord.Role, error)
	CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, error)
	DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error

	GetMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Member, error)
	GetMembers(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Member, error)
	SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...RequestOpt) ([]discord.Member, error)
	AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, error)
	RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error
	UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, error)
	AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error
	RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error

	UpdateSelfNick(guildID discord.Snowflake, nick string, opts ...RequestOpt) (*string, error)

	GetBans(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Ban, error)
	GetBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Ban, error)
	AddBan(guildID discord.Snowflake, userID discord.Snowflake, deleteMessageDays int, opts ...RequestOpt) error
	DeleteBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error

	GetIntegrations(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Integration, error)
	DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake, opts ...RequestOpt) error
	GetWebhooks(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Webhook, error)

	UpdateCurrentUserVoiceState(guildID discord.Snowflake, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
	UpdateUserVoiceState(guildID discord.Snowflake, userID discord.Snowflake, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
}

type guildServiceImpl struct {
	restClient Client
}

func (s *guildServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *guildServiceImpl) GetGuild(guildID discord.Snowflake, withCounts bool, opts ...RequestOpt) (guild *discord.Guild, err error) {
	values := route.QueryValues{}
	if withCounts {
		values["withCounts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuild.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guild, opts...)
	return
}

func (s *guildServiceImpl) GetGuildPreview(guildID discord.Snowflake, opts ...RequestOpt) (guildPreview *discord.GuildPreview, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildPreview.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildPreview, opts...)
	return
}

func (s *guildServiceImpl) CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (guild *discord.Guild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuild.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildCreate, &guild, opts...)
	return
}

func (s *guildServiceImpl) UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (guild *discord.Guild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuild.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildUpdate, &guild, opts...)
	return
}

func (s *guildServiceImpl) DeleteGuild(guildID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) GetRoles(guildID discord.Snowflake, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRoles.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &roles, opts...)
	return
}

func (s *guildServiceImpl) GetRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) (role []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &role, opts...)
	return
}

func (s *guildServiceImpl) CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateRole.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, createRole, &role, opts...)
	return
}

func (s *guildServiceImpl) UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, roleUpdate, &role, opts...)
	return
}

func (s *guildServiceImpl) UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRolePositions.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, rolePositionUpdates, &roles, opts...)
	return
}

func (s *guildServiceImpl) DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteRole.Compile(nil, guildID, roleID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) GetMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &member, opts...)
	return
}

func (s *guildServiceImpl) GetMembers(guildID discord.Snowflake, opts ...RequestOpt) (members []discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMembers.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *guildServiceImpl) SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...RequestOpt) (members []discord.Member, err error) {
	values := route.QueryValues{}
	if query != "" {
		values["query"] = query
	}
	if limit != 0 {
		values["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SearchMembers.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *guildServiceImpl) AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMembers.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, memberAdd, &member, opts...)
	return
}

func (s *guildServiceImpl) RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMember.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, memberUpdate, &member, opts...)
	return
}

func (s *guildServiceImpl) AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) UpdateSelfNick(guildID discord.Snowflake, nick string, opts ...RequestOpt) (nickName *string, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, discord.SelfNickUpdate{Nick: nick}, nickName, opts...)
	return
}

func (s *guildServiceImpl) GetBans(guildID discord.Snowflake, opts ...RequestOpt) (bans []discord.Ban, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBans.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &bans, opts...)
	return
}

func (s *guildServiceImpl) GetBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (ban *discord.Ban, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBan.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &ban, opts...)
	return
}

func (s *guildServiceImpl) AddBan(guildID discord.Snowflake, userID discord.Snowflake, deleteMessageDays int, opts ...RequestOpt) error {
	compiledRoute, err := route.AddBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, discord.AddBan{DeleteMessageDays: deleteMessageDays}, nil, opts...)
}

func (s *guildServiceImpl) DeleteBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) GetIntegrations(guildID discord.Snowflake, opts ...RequestOpt) (integrations []discord.Integration, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetIntegrations.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &integrations, opts...)
	return
}

func (s *guildServiceImpl) DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteIntegration.Compile(nil, guildID, integrationID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildServiceImpl) GetWebhooks(guildID discord.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildWebhooks.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *guildServiceImpl) GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *guildServiceImpl) UpdateCurrentUserVoiceState(guildID discord.Snowflake, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateCurrentUserVoiceState.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, currentUserVoiceStateUpdate, nil, opts...)
}

func (s *guildServiceImpl) UpdateUserVoiceState(guildID discord.Snowflake, userID discord.Snowflake, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateUserVoiceState.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, userVoiceStateUpdate, nil, opts...)
}
