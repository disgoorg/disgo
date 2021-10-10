package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ GuildTemplateService = (*guildTemplateServiceImpl)(nil)

func NewGuildTemplateService(restClient Client) GuildTemplateService {
	return &guildTemplateServiceImpl{restClient: restClient}
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
	GetGuildTemplates(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.GuildTemplate, error)
	CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (*discord.GuildTemplate, error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (*discord.Guild, error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (*discord.GuildTemplate, error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, error)
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

func (s *guildTemplateServiceImpl) GetGuildTemplates(guildID discord.Snowflake, opts ...RequestOpt) (guildTemplates []discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildTemplates.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplates, opts...)
	return
}

func (s *guildTemplateServiceImpl) CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
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

func (s *guildTemplateServiceImpl) SyncGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SyncGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, guildTemplateUpdate, &guildTemplate, opts...)
	return
}

func (s *guildTemplateServiceImpl) DeleteGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}
