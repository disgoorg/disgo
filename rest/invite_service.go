package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service = (*inviteServiceImpl)(nil)
	_ InviteService = (*inviteServiceImpl)(nil)
)

func NewInviteService(restClient Client) InviteService {
	return &inviteServiceImpl{restClient: restClient}
}

type InviteService interface {
	Service
	GetInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (*discord.Invite, error)
	DeleteInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	GetGuildInvites(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, error)
	GetChannelInvites(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, error)
}

type inviteServiceImpl struct {
	restClient Client
}

func (s *inviteServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *inviteServiceImpl) GetInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetInvite.Compile(nil, code)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteServiceImpl) CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateInvite.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, inviteCreate, &invite, opts...)
	return
}

func (s *inviteServiceImpl) DeleteInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteInvite.Compile(nil, code)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteServiceImpl) GetGuildInvites(guildID discord.Snowflake, opts ...RequestOpt) (invites []discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildInvites.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}

func (s *inviteServiceImpl) GetChannelInvites(channelID discord.Snowflake, opts ...RequestOpt) (invites []discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannelInvites.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}
