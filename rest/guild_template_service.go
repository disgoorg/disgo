package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewGuildTemplateService(client Client) GuildTemplateService {
	return nil
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}
