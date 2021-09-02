package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ GuildTemplateService = (*GuildTemplateServiceImpl)(nil)

func NewGuildTemplateService(restClient Client) GuildTemplateService {
	return &GuildTemplateServiceImpl{restClient: restClient}
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
}

type GuildTemplateServiceImpl struct {
	restClient Client
}

func (s *GuildTemplateServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *GuildTemplateServiceImpl) GetGuildTemplate(templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.GetGuildTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}

func (s *GuildTemplateServiceImpl) GetGuildTemplates(guildID discord.Snowflake, opts ...RequestOpt) (guildTemplates []discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.GetGuildTemplates.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guildTemplates, opts...)
	return
}

func (s *GuildTemplateServiceImpl) CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.CreateGuildTemplate.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, guildTemplateCreate, &guildTemplate, opts...)
	return
}

func (s *GuildTemplateServiceImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (guild *discord.Guild, rErr Error) {
	compiledRoute, err := route.CreateGuildFromTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, createGuildFromTemplate, &guild, opts...)
	return
}

func (s *GuildTemplateServiceImpl) SyncGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.SyncGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}

func (s *GuildTemplateServiceImpl) UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.UpdateGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, guildTemplateUpdate, &guildTemplate, opts...)
	return
}

func (s *GuildTemplateServiceImpl) DeleteGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (guildTemplate *discord.GuildTemplate, rErr Error) {
	compiledRoute, err := route.DeleteGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &guildTemplate, opts...)
	return
}
