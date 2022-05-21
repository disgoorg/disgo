package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Members = (*memberImpl)(nil)

func NewMembers(client Client) Members {
	return &memberImpl{client: client}
}

type Members interface {
	GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*discord.Member, error)
	GetMembers(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Member, error)
	SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) ([]discord.Member, error)
	AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, error)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, error)

	AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error
	RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error

	UpdateSelfNick(guildID snowflake.ID, nick string, opts ...RequestOpt) (*string, error)

	UpdateCurrentUserVoiceState(guildID snowflake.ID, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
	UpdateUserVoiceState(guildID snowflake.ID, userID snowflake.ID, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error
}

type memberImpl struct {
	client Client
}

func (s *memberImpl) GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) GetMembers(guildID snowflake.ID, opts ...RequestOpt) (members []discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMembers.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &members, opts...)
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) (members []discord.Member, err error) {
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
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd discord.MemberAdd, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.AddMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, memberAdd, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMember.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, memberUpdate, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *memberImpl) UpdateSelfNick(guildID snowflake.ID, nick string, opts ...RequestOpt) (nickName *string, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, discord.SelfNickUpdate{Nick: nick}, nickName, opts...)
	return
}

func (s *memberImpl) UpdateCurrentUserVoiceState(guildID snowflake.ID, currentUserVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateCurrentUserVoiceState.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, currentUserVoiceStateUpdate, nil, opts...)
}

func (s *memberImpl) UpdateUserVoiceState(guildID snowflake.ID, userID snowflake.ID, userVoiceStateUpdate discord.UserVoiceStateUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateUserVoiceState.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, userVoiceStateUpdate, nil, opts...)
}
