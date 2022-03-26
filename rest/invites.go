package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Invites = (*inviteImpl)(nil)

func NewInvites(restClient Client) Invites {
	return &inviteImpl{restClient: restClient}
}

type Invites interface {
	GetInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	CreateInvite(channelID snowflake.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (*discord.Invite, error)
	DeleteInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	GetGuildInvites(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Invite, error)
	GetChannelInvites(channelID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Invite, error)
}

type inviteImpl struct {
	restClient Client
}

func (s *inviteImpl) GetInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetInvite.Compile(nil, code)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteImpl) CreateInvite(channelID snowflake.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateInvite.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, inviteCreate, &invite, opts...)
	return
}

func (s *inviteImpl) DeleteInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteInvite.Compile(nil, code)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invite, opts...)
	return
}

func (s *inviteImpl) GetGuildInvites(guildID snowflake.Snowflake, opts ...RequestOpt) (invites []discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildInvites.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}

func (s *inviteImpl) GetChannelInvites(channelID snowflake.Snowflake, opts ...RequestOpt) (invites []discord.Invite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannelInvites.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &invites, opts...)
	return
}
