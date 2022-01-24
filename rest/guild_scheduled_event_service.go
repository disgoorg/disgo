package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Service                    = (*guildScheduledEventServiceImpl)(nil)
	_ GuildScheduledEventService = (*guildScheduledEventServiceImpl)(nil)
)

func NewGuildScheduledEventService(restClient Client) GuildScheduledEventService {
	return &guildScheduledEventServiceImpl{restClient: restClient}
}

type GuildScheduledEventService interface {
	Service
	GetGuildScheduledEvents(guildID snowflake.Snowflake, withUserCounts bool, opts ...RequestOpt) ([]discord.GuildScheduledEvent, error)
	GetGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, withUserCounts bool, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	CreateGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	UpdateGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...RequestOpt) (*discord.GuildScheduledEvent, error)
	DeleteGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, opts ...RequestOpt) error

	GetGuildScheduledEventUsers(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, limit int, withMember bool, before snowflake.Snowflake, after snowflake.Snowflake, opts ...RequestOpt) ([]discord.GuildScheduledEventUser, error)
}

type guildScheduledEventServiceImpl struct {
	restClient Client
}

func (s *guildScheduledEventServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *guildScheduledEventServiceImpl) GetGuildScheduledEvents(guildID snowflake.Snowflake, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvents []discord.GuildScheduledEvent, err error) {
	queryValues := route.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEvents.Compile(queryValues, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildScheduledEvents, opts...)
	return
}

func (s *guildScheduledEventServiceImpl) GetGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, withUserCounts bool, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	queryValues := route.QueryValues{}
	if withUserCounts {
		queryValues["with_user_counts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEvent.Compile(queryValues, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventServiceImpl) CreateGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildScheduledEvent.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildScheduledEventCreate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventServiceImpl) UpdateGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...RequestOpt) (guildScheduledEvent *discord.GuildScheduledEvent, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildScheduledEventUpdate, &guildScheduledEvent, opts...)
	return
}

func (s *guildScheduledEventServiceImpl) DeleteGuildScheduledEvent(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildScheduledEvent.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildScheduledEventServiceImpl) GetGuildScheduledEventUsers(guildID snowflake.Snowflake, guildScheduledEventID snowflake.Snowflake, limit int, withMember bool, before snowflake.Snowflake, after snowflake.Snowflake, opts ...RequestOpt) (guildScheduledEventUsers []discord.GuildScheduledEventUser, err error) {
	queryValues := route.QueryValues{}
	if limit > 0 {
		queryValues["limit"] = limit
	}
	if withMember {
		queryValues["withMember"] = true
	}
	if before != "" {
		queryValues["before"] = before
	}
	if after != "" {
		queryValues["after"] = after
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildScheduledEventUsers.Compile(nil, guildID, guildScheduledEventID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildScheduledEventUsers, opts...)
	return
}
