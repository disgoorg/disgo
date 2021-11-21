package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type GuildTemplate struct {
	discord.GuildTemplate
	Bot     *Bot
	Creator *User
}

// Guild returns the full Guild of the GuildTemplate if in caches
func (t *GuildTemplate) Guild() *Guild {
	return t.Bot.Caches.Guilds().Get(t.GuildID)
}

// Update updates the GuildTemplate with the provided UpdateGuildTemplate
func (t *GuildTemplate) Update(guildTemplateUpdate discord.GuildTemplateUpdate, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := t.Bot.RestServices.GuildTemplateService().UpdateGuildTemplate(t.GuildID, t.Code, guildTemplateUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return t.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// Sync updates the GuildTemplate with the provided UpdateGuildTemplate
func (t *GuildTemplate) Sync(opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := t.Bot.RestServices.GuildTemplateService().SyncGuildTemplate(t.GuildID, t.Code, opts...)
	if err != nil {
		return nil, err
	}
	return t.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// Delete deletes the GuildTemplate
func (t *GuildTemplate) Delete(opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := t.Bot.RestServices.GuildTemplateService().DeleteGuildTemplate(t.GuildID, t.Code, opts...)
	if err != nil {
		return nil, err
	}
	return t.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// CreateGuild creates a Guild from this GuildTemplate
func (t *GuildTemplate) CreateGuild(createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, error) {
	guild, err := t.Bot.RestServices.GuildTemplateService().CreateGuildFromTemplate(t.Code, createGuildFromTemplate, opts...)
	if err != nil {
		return nil, err
	}
	return t.Bot.EntityBuilder.CreateGuild(*guild, CacheStrategyNoWs), nil
}
