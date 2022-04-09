package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Threads = (*threadImpl)(nil)

func NewThreads(restClient Client) Threads {
	return &threadImpl{restClient: restClient}
}

type Threads interface {
	// CreateThreadFromMessage does not work for discord.ChannelTypeGuildForum channels.
	CreateThreadFromMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, threadCreateFromMessage discord.ThreadCreateFromMessage, opts ...RequestOpt) (thread discord.GuildThread, err error)
	CreateThreadInForum(channelID snowflake.Snowflake, threadCreateInForum discord.ThreadCreateInForum, opts ...RequestOpt) (thread discord.GuildThread, err error)
	CreateThread(channelID snowflake.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error)
	JoinThread(threadID snowflake.Snowflake, opts ...RequestOpt) error
	LeaveThread(threadID snowflake.Snowflake, opts ...RequestOpt) error
	AddThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error
	RemoveThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error
	GetThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error)
	GetThreadMembers(threadID snowflake.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error)

	GetPublicArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetPrivateArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetJoinedPrivateArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
}

type threadImpl struct {
	restClient Client
}

func (s *threadImpl) CreateThreadFromMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, threadCreateWithMessage discord.ThreadCreateFromMessage, opts ...RequestOpt) (thread discord.GuildThread, err error) {
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

func (s *threadImpl) CreateThreadInForum(channelID snowflake.Snowflake, threadCreateInForum discord.ThreadCreateInForum, opts ...RequestOpt) (thread discord.GuildThread, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateThread.Compile(nil, channelID)
	if err != nil {
		return
	}
	var channel discord.UnmarshalChannel
	err = s.restClient.Do(compiledRoute, threadCreateInForum, &channel, opts...)
	if err == nil {
		thread = channel.Channel.(discord.GuildThread)
	}
	return
}

func (s *threadImpl) CreateThread(channelID snowflake.Snowflake, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread discord.GuildThread, err error) {
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

func (s *threadImpl) JoinThread(threadID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.JoinThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) LeaveThread(threadID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.LeaveThread.Compile(nil, threadID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) AddThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.AddThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) RemoveThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *threadImpl) GetThreadMember(threadID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMember.Compile(nil, threadID, userID)
	if err != nil {
		return nil, err
	}
	err = s.restClient.Do(compiledRoute, nil, &threadMember, opts...)
	return
}

func (s *threadImpl) GetThreadMembers(threadID snowflake.Snowflake, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetThreadMembers.Compile(nil, threadID)
	if err != nil {
		return nil, err
	}
	err = s.restClient.Do(compiledRoute, nil, &threadMembers, opts...)
	return
}

func (s *threadImpl) GetPublicArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
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

func (s *threadImpl) GetPrivateArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
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

func (s *threadImpl) GetJoinedPrivateArchivedThreads(channelID snowflake.Snowflake, before discord.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
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
