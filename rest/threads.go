package rest

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Threads = (*threadImpl)(nil)

func NewThreads(client Client) Threads {
	return &threadImpl{client: client}
}

type Threads interface {
	CreateThreadWithMessage(channelID snowflake.ID, messageID snowflake.ID, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread discord.GuildThread, err error)
	CreateThread(channelID snowflake.ID, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error)
	JoinThread(threadID snowflake.ID, opts ...RequestOpt) error
	LeaveThread(threadID snowflake.ID, opts ...RequestOpt) error
	AddThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	GetThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error)
	GetThreadMembers(threadID snowflake.ID, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error)

	GetPublicArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetJoinedPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
}

type threadImpl struct {
	client Client
}

func (s *threadImpl) CreateThreadWithMessage(channelID snowflake.ID, messageID snowflake.ID, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...RequestOpt) (thread discord.GuildThread, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateThreadWithMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	var channel discord.UnmarshalChannel
	err = s.client.Do(compiledRoute, threadCreateWithMessage, &channel, opts...)
	if err == nil {
		thread = channel.Channel.(discord.GuildThread)
	}
	return
}

func (s *threadImpl) CreateThread(channelID snowflake.ID, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateThread.Compile(nil, channelID)
	if err != nil {
		return
	}
	var channel discord.UnmarshalChannel
	err = s.client.Do(compiledRoute, threadCreate, &channel, opts...)
	if err == nil {
		thread = channel.Channel.(discord.GuildThread)
	}
	return
}

func (s *threadImpl) JoinThread(threadID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.JoinThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) LeaveThread(threadID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.LeaveThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) AddThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.AddThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) GetThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(compiledRoute, nil, &threadMember, opts...)
	return
}

func (s *threadImpl) GetThreadMembers(threadID snowflake.ID, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMembers.Compile(nil, threadID)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(compiledRoute, nil, &threadMembers, opts...)
	return
}

func (s *threadImpl) GetPublicArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetArchivedPublicThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadImpl) GetPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetArchivedPrivateThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &threads, opts...)
	return
}

func (s *threadImpl) GetJoinedPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := route.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetJoinedAchievedPrivateThreads.Compile(queryValues, channelID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &threads, opts...)
	return
}
