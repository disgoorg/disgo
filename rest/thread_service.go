package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ ThreadService = (*threadServiceImpl)(nil)

func NewThreadService(restClient Client) ThreadService {
	return &threadServiceImpl{restClient: restClient}
}

type ThreadService interface {
	Service
	CreateThreadWithMessage(channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread *discord.Channel, rErr Error)
	CreateThread(channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread *discord.Channel, rErr Error)
	JoinThread(threadID discord.Snowflake, opts ...RequestOpt) Error
	LeaveThread(threadID discord.Snowflake, opts ...RequestOpt) Error
	AddThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error
	RemoveThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error
	GetThreadMembers(threadID discord.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, rErr Error)

	// GetActiveThreads will be removed in v10 in favor of GetActiveGuildThreads
	GetActiveThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error)
	GetArchivedPublicThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error)
	GetArchivedPrivateThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error)
	GetJoinedAchievedPrivateThreads(channelID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error)

	GetActiveGuildThreads(guildID discord.Snowflake, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error)
}

type threadServiceImpl struct {
	restClient Client
}

func (s *threadServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *threadServiceImpl) CreateThreadWithMessage(channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread *discord.Channel, rErr Error) {
	compiledRoute, err := route.CreateThreadWithMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, threadCreateWithMessage, &thread, opts...)
	return
}

func (s *threadServiceImpl) CreateThread(channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread *discord.Channel, rErr Error) {
	compiledRoute, err := route.CreateThread.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, threadCreate, &thread, opts...)
	return
}

func (s *threadServiceImpl) JoinThread(threadID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.JoinThread.Compile(nil, threadID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) LeaveThread(threadID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.LeaveThread.Compile(nil, threadID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) AddThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.AddThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) RemoveThreadMember(threadID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadServiceImpl) GetThreadMembers(threadID discord.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, rErr Error) {
	compiledRoute, err := route.GetThreadMembers.Compile(nil, threadID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, threadMembers, opts...)
	return
}

func (s *threadServiceImpl) GetActiveThreads(threadID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	compiledRoute, err := route.GetActiveThreads.Compile(queryValues, threadID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetArchivedPublicThreads(threadID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	compiledRoute, err := route.GetArchivedPublicThreads.Compile(queryValues, threadID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetArchivedPrivateThreads(threadID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	compiledRoute, err := route.GetArchivedPrivateThreads.Compile(queryValues, threadID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetJoinedAchievedPrivateThreads(threadID discord.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	compiledRoute, err := route.GetJoinedAchievedPrivateThreads.Compile(queryValues, threadID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadServiceImpl) GetActiveGuildThreads(guildID discord.Snowflake, opts ...RequestOpt) (threads *discord.GetThreads, rErr Error) {
	compiledRoute, err := route.GetActiveGuildThreads.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &threads, opts...)
	return
}
