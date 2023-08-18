package rest

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Threads = (*threadImpl)(nil)

func NewThreads(client Client) Threads {
	return &threadImpl{client: client}
}

type Threads interface {
	// CreateThreadFromMessage does not work for discord.ChannelTypeGuildForum or discord.ChannelTypeGuildMedia channels.
	CreateThreadFromMessage(channelID snowflake.ID, messageID snowflake.ID, threadCreateFromMessage discord.ThreadCreateFromMessage, opts ...RequestOpt) (thread *discord.GuildThread, err error)
	CreatePostInThreadChannel(channelID snowflake.ID, postCreateInChannel discord.ThreadChannelPostCreate, opts ...RequestOpt) (post *discord.ThreadChannelPost, err error)
	CreateThread(channelID snowflake.ID, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread *discord.GuildThread, err error)
	JoinThread(threadID snowflake.ID, opts ...RequestOpt) error
	LeaveThread(threadID snowflake.ID, opts ...RequestOpt) error
	AddThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	GetThreadMember(threadID snowflake.ID, userID snowflake.ID, withMember bool, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error)
	GetThreadMembers(threadID snowflake.ID, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error)
	GetThreadMembersPage(threadID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) ThreadMemberPage

	GetPublicArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
	GetJoinedPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error)
}

type threadImpl struct {
	client Client
}

func (s *threadImpl) CreateThreadFromMessage(channelID snowflake.ID, messageID snowflake.ID, threadCreateWithMessage discord.ThreadCreateFromMessage, opts ...RequestOpt) (thread *discord.GuildThread, err error) {
	err = s.client.Do(CreateThreadWithMessage.Compile(nil, channelID, messageID), threadCreateWithMessage, &thread, opts...)
	return
}

func (s *threadImpl) CreatePostInThreadChannel(channelID snowflake.ID, postCreateInChannel discord.ThreadChannelPostCreate, opts ...RequestOpt) (thread *discord.ThreadChannelPost, err error) {
	body, err := postCreateInChannel.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(CreateThread.Compile(nil, channelID), body, &thread, opts...)
	return
}

func (s *threadImpl) CreateThread(channelID snowflake.ID, threadCreate discord.ThreadCreate, opts ...RequestOpt) (thread *discord.GuildThread, err error) {
	err = s.client.Do(CreateThread.Compile(nil, channelID), threadCreate, &thread, opts...)
	return
}

func (s *threadImpl) JoinThread(threadID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(JoinThread.Compile(nil, threadID), nil, nil, opts...)
}

func (s *threadImpl) LeaveThread(threadID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(LeaveThread.Compile(nil, threadID), nil, nil, opts...)
}

func (s *threadImpl) AddThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(AddThreadMember.Compile(nil, threadID, userID), nil, nil, opts...)
}

func (s *threadImpl) RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveThreadMember.Compile(nil, threadID, userID), nil, nil, opts...)
}

func (s *threadImpl) GetThreadMember(threadID snowflake.ID, userID snowflake.ID, withMember bool, opts ...RequestOpt) (threadMember *discord.ThreadMember, err error) {
	err = s.client.Do(GetThreadMember.Compile(discord.QueryValues{"with_member": withMember}, threadID, userID), nil, &threadMember, opts...)
	return
}

func (s *threadImpl) GetThreadMembers(threadID snowflake.ID, opts ...RequestOpt) ([]discord.ThreadMember, error) {
	return s.getThreadMembers(threadID, nil, opts...)
}

func (s *threadImpl) GetThreadMembersPage(threadID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) ThreadMemberPage {
	return ThreadMemberPage{
		getItems: func(after snowflake.ID) ([]discord.ThreadMember, error) {
			queryValues := discord.QueryValues{
				"with_member": true,
				"after":       after,
			}
			if limit != 0 {
				queryValues["limit"] = limit
			}
			return s.getThreadMembers(threadID, queryValues, opts...)
		},
		ID: startID,
	}
}

func (s *threadImpl) GetPublicArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := discord.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	err = s.client.Do(GetPublicArchivedThreads.Compile(queryValues, channelID), nil, &threads, opts...)
	return
}

func (s *threadImpl) GetPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := discord.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	err = s.client.Do(GetPrivateArchivedThreads.Compile(queryValues, channelID), nil, &threads, opts...)
	return
}

func (s *threadImpl) GetJoinedPrivateArchivedThreads(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (threads *discord.GetThreads, err error) {
	queryValues := discord.QueryValues{}
	if !before.IsZero() {
		queryValues["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	err = s.client.Do(GetJoinedPrivateArchivedThreads.Compile(queryValues, channelID), nil, &threads, opts...)
	return
}

func (s *threadImpl) getThreadMembers(threadID snowflake.ID, queryValues discord.QueryValues, opts ...RequestOpt) (threadMembers []discord.ThreadMember, err error) {
	err = s.client.Do(GetThreadMembers.Compile(queryValues, threadID), nil, &threadMembers, opts...)
	return
}
