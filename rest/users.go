package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Users = (*userImpl)(nil)

func NewUsers(client Client) Users {
	return &userImpl{client: client}
}

type Users interface {
	GetUser(userID snowflake.ID, opts ...RequestOpt) (*discord.User, error)
	UpdateSelfUser(selfUserUpdate discord.SelfUserUpdate, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetGuilds(before int, after int, limit int, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error
	GetDMChannels(opts ...RequestOpt) ([]discord.Channel, error)
	CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (*discord.DMChannel, error)
}

type userImpl struct {
	client Client
}

func (s *userImpl) GetUser(userID snowflake.ID, opts ...RequestOpt) (user *discord.User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetUser.Compile(nil, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &user, opts...)
	return
}

func (s *userImpl) UpdateSelfUser(updateSelfUser discord.SelfUserUpdate, opts ...RequestOpt) (selfUser *discord.OAuth2User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateSelfUser.Compile(nil)
	if err != nil {
		return
	}
	var user *discord.User
	err = s.client.Do(compiledRoute, updateSelfUser, &user, opts...)
	return
}

func (s *userImpl) GetGuilds(before int, after int, limit int, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
	queryParams := route.QueryValues{}
	if before > 0 {
		queryParams["before"] = before
	}
	if after > 0 {
		queryParams["after"] = after
	}
	if limit > 0 {
		queryParams["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserGuilds.Compile(queryParams)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &guilds, opts...)
	return
}

func (s *userImpl) LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.LeaveGuild.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *userImpl) GetDMChannels(opts ...RequestOpt) (channels []discord.Channel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetDMChannels.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &channels, opts...)
	return
}

func (s *userImpl) CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (channel *discord.DMChannel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateDMChannel.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, discord.DMChannelCreate{RecipientID: userID}, &channel, opts...)
	return
}
