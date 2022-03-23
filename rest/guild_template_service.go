package rest

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
)

var (
	_ Service              = (*guildTemplateServiceImpl)(nil)
	_ GuildTemplateService = (*guildTemplateServiceImpl)(nil)
)

func NewGuildTemplateService(restClient Client) GuildTemplateService {
	return &guildTemplateServiceImpl{restClient: restClient}
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
	GetGuildTemplates(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.GuildTemplate, error)
	CreateGuildTemplate(guildID snowflake.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (*discord.GuildTemplate, error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (*discord.Guild, error)
	SyncGuildTemplate(guildID snowflake.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
	UpdateGuildTemplate(guildID snowflake.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (*discord.GuildTemplate, error)
	DeleteGuildTemplate(guildID snowflake.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
}

type guildTemplateServiceImpl struct {
	restClient Client
}

func (s *guildTemplateServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *guildTemplateServiceImpl) GetGuildTemplate(templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildTemplate.Compile(nil, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) GetGuildTemplates(guildID snowflake.Snowflake, opts ...RequestOpt) (guildTemplates []discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildTemplates.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplates, opts...)
	return
}

func (s *guildTemplateServiceImpl) CreateGuildTemplate(guildID snowflake.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildTemplate.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildTemplateCreate, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (guild *discord.Guild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildFromTemplate.Compile(nil, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, createGuildFromTemplate, &guild, opts...)
	return
}

func (s *guildTemplateServiceImpl) SyncGuildTemplate(guildID snowflake.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SyncGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) UpdateGuildTemplate(guildID snowflake.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildTemplateUpdate, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) DeleteGuildTemplate(guildID snowflake.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}
