package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ GuildScheduledEvents = (*guildScheduledEventImpl)(nil)

func NewGuildScheduledEvents(client Client) GuildScheduledEvents {
	return &guildScheduledEventImpl{client: client}
}

type GuildScheduledEvents interface {
	GetGuildScheduledEvents(guildID snowflake.ID, withUserCounts bool, opts ...RequestOpt) ([]discord.GuildScheduledEvent, error)
	GetGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withUserCounts bool, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	CreateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	UpdateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	DeleteGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) error

	GetGuildScheduledEventUsers(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withMember bool, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]discord.GuildScheduledEventUser, error)
	GetGuildScheduledEventUsersPage(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withMember bool, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.GuildScheduledEventUser]
}

type guildScheduledEventImpl struct {
	client Client
}

func (s *guildScheduledEventImpl) GetGuildScheduledEvents(guildID snowflake.ID, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvents []discord.GuildScheduledEvent, err error) {
	queryValues := discord.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	err = s.client.Do(GetGuildScheduledEvents.Compile(queryValues, guildID), nil, &guildScheduledEvents, opts...)
	return
}

func (s *guildScheduledEventImpl) GetGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	queryValues := discord.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	err = s.client.Do(GetGuildScheduledEvent.Compile(queryValues, guildID, guildScheduledEventID), nil, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) CreateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	err = s.client.Do(CreateGuildScheduledEvent.Compile(nil, guildID), guildScheduledEventCreate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) UpdateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	err = s.client.Do(UpdateGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID), guildScheduledEventUpdate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) DeleteGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID), nil, nil, opts...)
}

func (s *guildScheduledEventImpl) GetGuildScheduledEventUsers(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withMember bool, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (guildScheduledEventUsers []discord.GuildScheduledEventUser, err error) {
	queryValues := discord.QueryValues{}
	if limit > 0 {
		queryValues["limit"] = limit
	}
	if withMember {
		queryValues["withMember"] = true
	}
	if before != 0 {
		queryValues["before"] = before
	}
	if after != 0 {
		queryValues["after"] = after
	}
	err = s.client.Do(GetGuildScheduledEventUsers.Compile(nil, guildID, guildScheduledEventID), nil, &guildScheduledEventUsers, opts...)
	return
}

func (s *guildScheduledEventImpl) GetGuildScheduledEventUsersPage(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withMember bool, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.GuildScheduledEventUser] {
	return Page[discord.GuildScheduledEventUser]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) ([]discord.GuildScheduledEventUser, error) {
			return s.GetGuildScheduledEventUsers(guildID, guildScheduledEventID, withMember, before, after, limit, opts...)
		},
		getIDFunc: func(user discord.GuildScheduledEventUser) snowflake.ID {
			return user.User.ID
		},
		ID: startID,
	}
}
