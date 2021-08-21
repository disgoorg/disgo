package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewApplicationService(client Client) ApplicationService {
	return nil
}

type ApplicationService interface {
	Service
	GetGlobalCommands(ctx context.Context, applicationID discord.Snowflake) ([]discord.ApplicationCommand, Error)
	GetGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake) (*discord.ApplicationCommand, Error)
	CreateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, command discord.ApplicationCommandCreate) (*discord.ApplicationCommand, Error)
	SetGlobalCommands(ctx context.Context, applicationID discord.Snowflake, commands ...discord.ApplicationCommandCreate) ([]discord.ApplicationCommand, Error)
	UpdateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate) (*discord.ApplicationCommand, Error)
	DeleteGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake) Error

	GetGuildCommands(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake) ([]discord.ApplicationCommand, Error)
	GetGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (*discord.ApplicationCommand, Error)
	CreateGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, command discord.ApplicationCommandCreate) (*discord.ApplicationCommand, Error)
	SetGuildCommands(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commands ...discord.ApplicationCommandCreate) ([]discord.ApplicationCommand, Error)
	UpdateGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate) (*discord.ApplicationCommand, Error)
	DeleteGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) Error

	GetGuildCommandsPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake) ([]discord.GuildCommandPermissions, Error)
	GetGuildCommandPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (*discord.GuildCommandPermissions, Error)
	SetGuildCommandsPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) ([]discord.GuildCommandPermissions, Error)
	SetGuildCommandPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions ...discord.CommandPermission) (*discord.GuildCommandPermissions, Error)
}
