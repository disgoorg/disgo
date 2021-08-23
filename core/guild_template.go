package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type GuildTemplate struct {
	discord.GuildTemplate
	Disgo   Disgo
	Creator *User
}

// Guild returns the full Guild of the GuildTemplate if in cache
func (t *GuildTemplate) Guild() *Guild {
	return t.Disgo.Cache().GuildCache().Get(t.GuildID)
}

// Update updates the GuildTemplate with the provided UpdateGuildTemplate
func (t *GuildTemplate) Update(guildTemplateUpdate discord.GuildTemplateUpdate) (*GuildTemplate, rest.Error) {
	guildTemplate, err := t.Disgo.RestServices().GuildTemplateService().UpdateGuildTemplate(t.GuildID, t.Code, guildTemplateUpdate)
	if err != nil {
		return nil, err
	}
	return t.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// Sync updates the GuildTemplate with the provided UpdateGuildTemplate
func (t *GuildTemplate) Sync(opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := t.Disgo.RestServices().GuildTemplateService().SyncGuildTemplate(t.GuildID, t.Code)
	if err != nil {
		return nil, err
	}
	return t.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// Delete deletes the GuildTemplate
func (t *GuildTemplate) Delete(opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := t.Disgo.RestServices().GuildTemplateService().DeleteGuildTemplate(t.GuildID, t.Code)
	if err != nil {
		return nil, err
	}
	return t.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}
