package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service = (*threadServiceImpl)(nil)
	_ ThreadService = (*threadServiceImpl)(nil)
)

func NewThreadService(restClient Client) ThreadService {
	return &threadServiceImpl{restClient: restClient}
}

type ThreadService interface {
	Service
	CreateThreadWithMessage(channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread discord.GuildThread, err error)
	CreateThread(channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error)
	JoinThread(threadID discord.Snowflake, opts ...RequestOpt) error
	LeaveThread(threadID discord.Snowflake, opts ...RequestOpt) error
	AddThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error
	RemoveThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error
	GetThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error)
	GetThreadMembers(threadID discord.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error)

	GetPublicArchivedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetPrivateArchivedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetJoinedPrivateAchievedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)

	GetActiveGuildThreads(guildID discord.Snowflake, opts ...RequestOpt) (threads *discord.GetAllThreads, err error)
}

type threadServiceImpl struct {
	restClient Client
}

func (s *threadServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *threadServiceImpl) CreateThreadWithMessage(channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread discord.GuildThread, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateThreadWithMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	var channel discord.UnmarshalChannel
	err = s.restClient.Do(compiledRoute, threadCreateWithMessage, &channel, opts...)
	if err == nil {
		thread = channel.Channel.(discord.GuildThread)
	}
	return
}

func (s *threadServiceImpl) CreateThread(channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateThread.Compile(nil, channelID)
	if err != nil {
		return
	}
	var channel discord.UnmarshalChannel
	err = s.restClient.Do(compiledRoute, threadCreate, &channel, opts...)
	if err == nil {
		thread = channel.Channel.(discord.GuildThread)
	}
	return
}

func (s *threadServiceImpl) JoinThread(threadID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.JoinThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) LeaveThread(threadID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.LeaveThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) AddThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.AddThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) RemoveThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) GetThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return nil, err
	}
	err = s.restClient.Do(compiledRoute, nil, &threadMember, opts...)
	return
}

func (s *threadServiceImpl) GetThreadMembers(threadID discord.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMembers.Compile(nil, threadID)
	if err != nil {
		return nil, err
	}
	err = s.restClient.Do(compiledRoute, nil, &threadMembers, opts...)
	return
}

func (s *threadServiceImpl) GetPublicArchivedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetArchivedPublicThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetPrivateArchivedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetArchivedPrivateThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetJoinedPrivateAchievedThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetJoinedAchievedPrivateThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetActiveGuildThreads(guildID discord.Snowflake, opts ...RequestOpt) (threads *discord.GetAllThreads, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetActiveGuildThreads.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}
