package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Members = (*memberImpl)(nil)

func NewMembers(client Client) Members {
	return &memberImpl{client: client}
}

type Members interface {
	GetMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (*discord.Member, error)
	GetMembers(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Member, error)
	SearchMembers(guildID snowflake.Snowflake, query string, limit int, opts ...RequestOpt) ([]discord.Member, error)
	AddMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, error)
	RemoveMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error
	UpdateMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, error)

	AddMemberRole(guildID snowflake.Snowflake, userID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error
	RemoveMemberRole(guildID snowflake.Snowflake, userID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error

	UpdateSelfNick(guildID snowflake.Snowflake, nick string, opts ...RequestOpt) (*string, error)

	UpdateCurrentUserVoiceState(guildID snowflake.Snowflake, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
	UpdateUserVoiceState(guildID snowflake.Snowflake, userID snowflake.Snowflake, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
}

type memberImpl struct {
	client Client
}

func (s *memberImpl) GetMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &member, opts...)
	return
}

func (s *memberImpl) GetMembers(guildID snowflake.Snowflake, opts ...RequestOpt) (members []discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMembers.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *memberImpl) SearchMembers(guildID snowflake.Snowflake, query string, limit int, opts ...RequestOpt) (members []discord.Member, err error) {
	values := route.QueryValues{}
	if query != "" {
		values["query"] = query
	}
	if limit != 0 {
		values["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SearchMembers.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &members, opts...)
	return
}

func (s *memberImpl) AddMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.AddMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, memberAdd, &member, opts...)
	return
}

func (s *memberImpl) RemoveMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMember.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) UpdateMember(guildID snowflake.Snowflake, userID snowflake.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, memberUpdate, &member, opts...)
	return
}

func (s *memberImpl) AddMemberRole(guildID snowflake.Snowflake, userID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) RemoveMemberRole(guildID snowflake.Snowflake, userID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) UpdateSelfNick(guildID snowflake.Snowflake, nick string, opts ...RequestOpt) (nickName *string, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, discord.SelfNickUpdate{Nick: nick}, nickName, opts...)
	return
}

func (s *memberImpl) UpdateCurrentUserVoiceState(guildID snowflake.Snowflake, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateCurrentUserVoiceState.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, currentUserVoiceStateUpdate, nil, opts...)
}

func (s *memberImpl) UpdateUserVoiceState(guildID snowflake.Snowflake, userID snowflake.Snowflake, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateUserVoiceState.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, userVoiceStateUpdate, nil, opts...)
}
