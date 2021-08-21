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
	GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, Error)
	GetSelfUser(ctx context.Context) (*discord.SelfUser, Error)
	UpdateSelfUser(ctx context.Context, updateSelfUser discord.UpdateSelfUser) (*discord.SelfUser, Error)
	GetGuilds(ctx context.Context, before int, after int, limit int) ([]discord.PartialGuild, Error)
	LeaveGuild(ctx context.Context, guildID discord.Snowflake) Error
	GetDMChannels(ctx context.Context) ([]discord.Channel, Error)
	CreateDMChannel(ctx context.Context, userID discord.Snowflake) (*discord.Channel, Error)
}

type UserServiceImpl struct {
	restClient Client
}

func (s *UserServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *UserServiceImpl) GetUser(ctx context.Context, userID discord.Snowflake) (user *discord.User, rErr Error) {
	compiledRoute, err := route.GetUser.Compile(nil, userID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &user)
	return
}

func (s *UserServiceImpl) GetSelfUser(ctx context.Context) (selfUser *discord.SelfUser, rErr Error) {
	compiledRoute, err := route.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &selfUser)
	return
}

func (s *UserServiceImpl) UpdateSelfUser(ctx context.Context, updateSelfUser discord.UpdateSelfUser) (selfUser *discord.SelfUser, rErr Error) {
	compiledRoute, err := route.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	var user *discord.User
	rErr = s.restClient.Do(ctx, compiledRoute, updateSelfUser, &user)
	return
}

func (s *UserServiceImpl) GetGuilds(ctx context.Context, before int, after int, limit int) (guilds []discord.PartialGuild, rErr Error) {
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

	rErr = s.restClient.Do(ctx, compiledRoute, nil, &guilds)
	return
}

func (s *UserServiceImpl) LeaveGuild(ctx context.Context, guildID discord.Snowflake) Error {
	compiledRoute, err := route.LeaveGuild.Compile(nil, guildID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(ctx, compiledRoute, nil, nil)
}

func (s *UserServiceImpl) GetDMChannels(ctx context.Context) (channels []discord.Channel, rErr Error) {
	compiledRoute, err := route.GetDMChannels.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(ctx, compiledRoute, nil, &channels)
	return
}

func (s *UserServiceImpl) CreateDMChannel(ctx context.Context, userID discord.Snowflake) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.CreateDMChannel.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(ctx, compiledRoute, discord.DMChannelCreate{RecipientID: userID}, &channel)
	return
}
