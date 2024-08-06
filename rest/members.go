package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Members = (*memberImpl)(nil)

func NewMembers(client Client) Members {
	return &memberImpl{client: client}
}

type Members interface {
	GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*discord.Member, error)
	GetMembers(guildID snowflake.ID, limit int, after snowflake.ID, opts ...RequestOpt) ([]discord.Member, error)
	SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) ([]discord.Member, error)
	AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, error)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, error)

	AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error
	RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error

	UpdateCurrentMember(guildID snowflake.ID, nick string, opts ...RequestOpt) (*string, error)

	GetCurrentUserVoiceState(guildID snowflake.ID, opts ...RequestOpt) (*discord.VoiceState, error)
	GetUserVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*discord.VoiceState, error)
	UpdateCurrentUserVoiceState(guildID snowflake.ID, currentUserVoiceStateUpdate discord.CurrentUserVoiceStateUpdate, opts ...RequestOpt) error
	UpdateUserVoiceState(guildID snowflake.ID, userID snowflake.ID, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
}

type memberImpl struct {
	client Client
}

func (s *memberImpl) GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (member *discord.Member, err error) {
	err = s.client.Do(GetMember.Compile(nil, guildID, userID), nil, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) GetMembers(guildID snowflake.ID, limit int, after snowflake.ID, opts ...RequestOpt) (members []discord.Member, err error) {
	values := discord.QueryValues{
		"limit": limit,
		"after": after,
	}
	err = s.client.Do(GetMembers.Compile(values, guildID), nil, &members, opts...)
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) (members []discord.Member, err error) {
	values := discord.QueryValues{}
	if query != "" {
		values["query"] = query
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(SearchMembers.Compile(values, guildID), nil, &members, opts...)
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd discord.MemberAdd, opts ...RequestOpt) (member *discord.Member, err error) {
	err = s.client.Do(AddMember.Compile(nil, guildID, userID), memberAdd, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveMember.Compile(nil, guildID, userID), nil, nil, opts...)
}

func (s *memberImpl) UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (member *discord.Member, err error) {
	err = s.client.Do(UpdateMember.Compile(nil, guildID, userID), memberUpdate, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(AddMemberRole.Compile(nil, guildID, userID, roleID), nil, nil, opts...)
}

func (s *memberImpl) RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveMemberRole.Compile(nil, guildID, userID, roleID), nil, nil, opts...)
}

func (s *memberImpl) UpdateCurrentMember(guildID snowflake.ID, nick string, opts ...RequestOpt) (nickName *string, err error) {
	err = s.client.Do(UpdateCurrentMember.Compile(nil, guildID), discord.CurrentMemberUpdate{Nick: nick}, nickName, opts...)
	return
}

func (s *memberImpl) GetCurrentUserVoiceState(guildID snowflake.ID, opts ...RequestOpt) (state *discord.VoiceState, err error) {
	err = s.client.Do(GetCurrentUserVoiceState.Compile(nil, guildID), nil, &state, opts...)
	return
}

func (s *memberImpl) GetUserVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (state *discord.VoiceState, err error) {
	err = s.client.Do(GetUserVoiceState.Compile(nil, guildID, userID), nil, &state, opts...)
	return
}

func (s *memberImpl) UpdateCurrentUserVoiceState(guildID snowflake.ID, currentUserVoiceStateUpdate discord.CurrentUserVoiceStateUpdate, opts ...RequestOpt) error {
	return s.client.Do(UpdateCurrentUserVoiceState.Compile(nil, guildID), currentUserVoiceStateUpdate, nil, opts...)
}

func (s *memberImpl) UpdateUserVoiceState(guildID snowflake.ID, userID snowflake.ID, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	return s.client.Do(UpdateUserVoiceState.Compile(nil, guildID, userID), userVoiceStateUpdate, nil, opts...)
}
