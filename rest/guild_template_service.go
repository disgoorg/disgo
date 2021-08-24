package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewGuildTemplateService(client Client) GuildTemplateService {
	return nil
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
