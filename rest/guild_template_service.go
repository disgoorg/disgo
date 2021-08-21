package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewGuildTemplateService(client Client) GuildTemplateService {
	return nil
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(ctx context.Context, templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(ctx context.Context, guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(ctx context.Context, guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(ctx context.Context, templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}
