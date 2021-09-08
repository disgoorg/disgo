package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ InviteService = (*inviteServiceImpl)(nil)

func NewInviteService(restClient Client) InviteService {
	return &inviteServiceImpl{restClient: restClient}
}

type InviteService interface {
	Service
	GetInvite(code string, opts ...RequestOpt) (*discord.Invite, Error)
	CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (*discord.Invite, Error)
	DeleteInvite(code string, opts ...RequestOpt) (*discord.Invite, Error)
	GetGuildInvites(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, Error)
	GetChannelInvites(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, Error)
}

type inviteServiceImpl struct {
	restClient Client
}

func (s *inviteServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *inviteServiceImpl) GetInvite(code string, opts ...RequestOpt) (invite *discord.Invite, rErr Error) {
	compiledRoute, err := route.GetInvite.Compile(nil, code)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteServiceImpl) CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (invite *discord.Invite, rErr Error) {
	compiledRoute, err := route.CreateInvite.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, inviteCreate, &invite, opts...)
	return
}

func (s *inviteServiceImpl) DeleteInvite(code string, opts ...RequestOpt) (invite *discord.Invite, rErr Error) {
	compiledRoute, err := route.DeleteInvite.Compile(nil, code)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteServiceImpl) GetGuildInvites(guildID discord.Snowflake, opts ...RequestOpt) (invites []discord.Invite, rErr Error) {
	compiledRoute, err := route.GetGuildInvites.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}

func (s *inviteServiceImpl) GetChannelInvites(channelID discord.Snowflake, opts ...RequestOpt) (invites []discord.Invite, rErr Error) {
	compiledRoute, err := route.GetChannelInvites.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}
