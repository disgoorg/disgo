package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Invites = (*inviteImpl)(nil)

func NewInvites(client Client) Invites {
	return &inviteImpl{client: client}
}

type Invites interface {
	GetInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	CreateInvite(channelID snowflake.ID, inviteCreate discord.InviteCreate, opts ...RequestOpt) (*discord.Invite, error)
	DeleteInvite(code string, opts ...RequestOpt) (*discord.Invite, error)
	GetGuildInvites(guildID snowflake.ID, opts ...RequestOpt) ([]discord.ExtendedInvite, error)
	GetChannelInvites(channelID snowflake.ID, opts ...RequestOpt) ([]discord.ExtendedInvite, error)
}

type inviteImpl struct {
	client Client
}

func (s *inviteImpl) GetInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	err = s.client.Do(GetInvite.Compile(nil, code), nil, &invite, opts...)
	return
}

func (s *inviteImpl) CreateInvite(channelID snowflake.ID, inviteCreate discord.InviteCreate, opts ...RequestOpt) (invite *discord.Invite, err error) {
	body, err := inviteCreate.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(CreateInvite.Compile(nil, channelID), body, &invite, opts...)
	return
}

func (s *inviteImpl) DeleteInvite(code string, opts ...RequestOpt) (invite *discord.Invite, err error) {
	err = s.client.Do(DeleteInvite.Compile(nil, code), nil, &invite, opts...)
	return
}

func (s *inviteImpl) GetInviteTargetUsers(code string, opts ...RequestOpt) (targetUsers *string, err error) {
	err = s.client.Do(GetInviteTargetUsers.Compile(nil, code), nil, &targetUsers, opts...)
	return
}

func (s *inviteImpl) SetInviteTargetUsers(code string, inviteTargetUsersUpdate discord.InviteTargetUsersUpdate, opts ...RequestOpt) (err error) {
	body, err := inviteTargetUsersUpdate.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(SetInviteTargetUsers.Compile(nil, code), body, nil, opts...)
	return
}

func (s *inviteImpl) GetInviteTargetUsersJobStatus(code string, opts ...RequestOpt) (targetUsersJobStatus *discord.TargetUsersJobStatus, err error) {
	err = s.client.Do(GetInviteTargetUsersJobStatus.Compile(nil, code), nil, &targetUsersJobStatus, opts...)
	return
}

func (s *inviteImpl) GetGuildInvites(guildID snowflake.ID, opts ...RequestOpt) (invites []discord.ExtendedInvite, err error) {
	err = s.client.Do(GetGuildInvites.Compile(nil, guildID), nil, &invites, opts...)
	return
}

func (s *inviteImpl) GetChannelInvites(channelID snowflake.ID, opts ...RequestOpt) (invites []discord.ExtendedInvite, err error) {
	err = s.client.Do(GetChannelInvites.Compile(nil, channelID), nil, &invites, opts...)
	return
}
