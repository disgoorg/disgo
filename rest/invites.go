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
	GetInviteTargetUsers(code string, opts ...RequestOpt) (targetUsers *string, err error)
	UpdateInviteTargetUsers(code string, inviteTargetUsersUpdate discord.InviteTargetUsersUpdate, opts ...RequestOpt) (err error)
	AddInviteTargetUser(code string, userID snowflake.ID, opts ...RequestOpt) error
	RemoveInviteTargetUser(code string, userID snowflake.ID, opts ...RequestOpt) error
	BulkAddInviteTargetUsers(code string, userIDs []snowflake.ID, opts ...RequestOpt) error
	BulkDeleteInviteTargetUsers(code string, userIDs []snowflake.ID, opts ...RequestOpt) error
	GetInviteTargetUsersJobStatus(code string, opts ...RequestOpt) (targetUsersJobStatus *discord.TargetUsersJobStatus, err error)
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
	err = s.client.Do(CreateInvite.Compile(nil, channelID), inviteCreate, &invite, opts...)
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

func (s *inviteImpl) UpdateInviteTargetUsers(code string, inviteTargetUsersUpdate discord.InviteTargetUsersUpdate, opts ...RequestOpt) (err error) {
	body, err := inviteTargetUsersUpdate.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(SetInviteTargetUsers.Compile(nil, code), body, nil, opts...)
	return
}

func (s *inviteImpl) AddInviteTargetUser(code string, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(AddInviteTargetUser.Compile(nil, code, userID), nil, nil, opts...)
}

func (s *inviteImpl) RemoveInviteTargetUser(code string, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveInviteTargetUser.Compile(nil, code, userID), nil, nil, opts...)
}

func (s *inviteImpl) BulkAddInviteTargetUsers(code string, userIDs []snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(BulkAddInviteTargetUsers.Compile(nil, code), discord.InviteTargetUsersBulk{UserIDs: userIDs}, nil, opts...)
}

func (s *inviteImpl) BulkDeleteInviteTargetUsers(code string, userIDs []snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(BulkDeleteInviteTargetUsers.Compile(nil, code), discord.InviteTargetUsersBulk{UserIDs: userIDs}, nil, opts...)
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
