package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ UserService = (*UserServiceImpl)(nil)

func NewUserService(restClient Client) UserService {
	return &UserServiceImpl{
		restClient: restClient,
	}
}

type UserService interface {
	Service
	GetUser(userID discord.Snowflake) (*discord.User, Error)
	GetSelfUser(opts ...rest.RequestOpt) (*discord.SelfUser, Error)
	UpdateSelfUser(updateSelfUser discord.UpdateSelfUser) (*discord.SelfUser, Error)
	GetGuilds(before int, after int, limit int) ([]discord.PartialGuild, Error)
	LeaveGuild(guildID discord.Snowflake) Error
	GetDMChannels(opts ...rest.RequestOpt) ([]discord.Channel, Error)
	CreateDMChannel(userID discord.Snowflake) (*discord.Channel, Error)
}

type UserServiceImpl struct {
	restClient Client
}

func (s *UserServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *UserServiceImpl) GetUser(userID discord.Snowflake) (user *discord.User, rErr Error) {
	compiledRoute, err := route.GetUser.Compile(nil, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &user)
	return
}

func (s *UserServiceImpl) GetSelfUser(opts ...rest.RequestOpt) (selfUser *discord.SelfUser, rErr Error) {
	compiledRoute, err := route.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &selfUser)
	return
}

func (s *UserServiceImpl) UpdateSelfUser(updateSelfUser discord.UpdateSelfUser) (selfUser *discord.SelfUser, rErr Error) {
	compiledRoute, err := route.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	var user *discord.User
	rErr = s.restClient.Do(compiledRoute, updateSelfUser, &user)
	return
}

func (s *UserServiceImpl) GetGuilds(before int, after int, limit int) (guilds []discord.PartialGuild, rErr Error) {
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
	compiledRoute, err := route.GetGuilds.Compile(queryParams)
	if err != nil {
		return nil, NewError(nil, NewError(nil, err))
	}

	rErr = s.restClient.Do(compiledRoute, nil, &guilds)
	return
}

func (s *UserServiceImpl) LeaveGuild(guildID discord.Snowflake) Error {
	compiledRoute, err := route.LeaveGuild.Compile(nil, guildID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil)
}

func (s *UserServiceImpl) GetDMChannels(opts ...rest.RequestOpt) (channels []discord.Channel, rErr Error) {
	compiledRoute, err := route.GetDMChannels.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, nil, &channels)
	return
}

func (s *UserServiceImpl) CreateDMChannel(userID discord.Snowflake) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.CreateDMChannel.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, discord.DMChannelCreate{RecipientID: userID}, &channel)
	return
}
