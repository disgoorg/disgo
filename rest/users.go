package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Users = (*userImpl)(nil)

func NewUsers(client Client) Users {
	return &userImpl{client: client}
}

type Users interface {
	GetUser(userID snowflake.ID, opts ...RequestOpt) (*discord.User, error)
	UpdateCurrentUser(userUpdate discord.UserUpdate, opts ...RequestOpt) (*discord.OAuth2User, error)
	LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error
	GetDMChannels(opts ...RequestOpt) ([]discord.Channel, error)
	CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (*discord.DMChannel, error)
}

type userImpl struct {
	client Client
}

func (s *userImpl) GetUser(userID snowflake.ID, opts ...RequestOpt) (user *discord.User, err error) {
	err = s.client.Do(GetUser.Compile(nil, userID), nil, &user, opts...)
	return
}

func (s *userImpl) UpdateCurrentUser(userUpdate discord.UserUpdate, opts ...RequestOpt) (selfUser *discord.OAuth2User, err error) {
	err = s.client.Do(UpdateCurrentUser.Compile(nil), userUpdate, &selfUser, opts...)
	return
}

func (s *userImpl) LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(LeaveGuild.Compile(nil, guildID), nil, nil, opts...)
}

func (s *userImpl) GetDMChannels(opts ...RequestOpt) (channels []discord.Channel, err error) {
	err = s.client.Do(GetDMChannels.Compile(nil), nil, &channels, opts...)
	return
}

func (s *userImpl) CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (channel *discord.DMChannel, err error) {
	err = s.client.Do(CreateDMChannel.Compile(nil), discord.DMChannelCreate{RecipientID: userID}, &channel, opts...)
	return
}
