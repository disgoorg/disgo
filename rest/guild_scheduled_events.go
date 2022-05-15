package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
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

	GetGuildScheduledEventUsers(guildID snowflake.ID, guildScheduledEventID snowflake.ID, limit int, withMember bool, before snowflake.ID, after snowflake.ID, opts ...RequestOpt) ([]discord.GuildScheduledEventUser, error)
}

type guildScheduledEventImpl struct {
	client Client
}

func (s *guildScheduledEventImpl) GetGuildScheduledEvents(guildID snowflake.ID, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvents []discord.GuildScheduledEvent, err error) {
	queryValues := route.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEvents.Compile(queryValues, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &guildScheduledEvents, opts...)
	return
}

func (s *guildScheduledEventImpl) GetGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	queryValues := route.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEvent.Compile(queryValues, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) CreateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildScheduledEvent.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, guildScheduledEventCreate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) UpdateGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, guildScheduledEventUpdate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventImpl) DeleteGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildScheduledEventImpl) GetGuildScheduledEventUsers(guildID snowflake.ID, guildScheduledEventID snowflake.ID, limit int, withMember bool, before snowflake.ID, after snowflake.ID, opts ...RequestOpt) (guildScheduledEventUsers []discord.GuildScheduledEventUser, err error) {
	queryValues := route.QueryValues{}
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

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEventUsers.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &guildScheduledEventUsers, opts...)
	return
}
