package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ ApplicationService = (*ApplicationServiceImpl)(nil)

func NewApplicationService(client Client) ApplicationService {
	return &ApplicationServiceImpl{restClient: client}
}

type ApplicationService interface {
	Service
	GetBotApplicationInfo(ctx context.Context) (*discord.Application, Error)
	GetAuthorizationInfo(ctx context.Context) (*discord.AuthorizationInformation, Error)

	GetGlobalCommands(ctx context.Context, applicationID discord.Snowflake) ([]discord.ApplicationCommand, Error)
	GetGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake) (*discord.ApplicationCommand, Error)
	CreateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate) (*discord.ApplicationCommand, Error)
	SetGlobalCommands(ctx context.Context, applicationID discord.Snowflake, commandCreates ...discord.ApplicationCommandCreate) ([]discord.ApplicationCommand, Error)
	UpdateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (*discord.ApplicationCommand, Error)
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

type ApplicationServiceImpl struct {
	restClient Client
}

func (s *ApplicationServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *ApplicationServiceImpl) GetBotApplicationInfo(ctx context.Context) (application *discord.Application, rErr Error) {
	compiledRoute, err := route.GetBotApplicationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &application)
	return
}

func (s *ApplicationServiceImpl) GetAuthorizationInfo(ctx context.Context) (info *discord.AuthorizationInformation, rErr Error) {
	compiledRoute, err := route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &info)
	return
}

func (s *ApplicationServiceImpl) GetGlobalCommands(ctx context.Context, applicationID discord.Snowflake) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commands)
	return
}

func (s *ApplicationServiceImpl) GetGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) CreateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.CreateGlobalCommand.Compile(nil, applicationID, commandCreate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) SetGlobalCommands(ctx context.Context, applicationID discord.Snowflake, commandCreates ...discord.ApplicationCommandCreate) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.SetGlobalCommands.Compile(nil, applicationID, commandCreates)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commands)
	return
}

func (s *ApplicationServiceImpl) UpdateGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.UpdateGlobalCommand.Compile(nil, applicationID, commandID, commandUpdate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) DeleteGlobalCommand(ctx context.Context, applicationID discord.Snowflake, commandID discord.Snowflake) Error {
	compiledRoute, err := route.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(ctx, compiledRoute, nil, nil)
}

func (s *ApplicationServiceImpl) GetGuildCommands(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake) (command []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) GetGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) CreateGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.CreateGuildCommand.Compile(nil, applicationID, guildID, commandCreate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommands(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandCreates ...discord.ApplicationCommandCreate) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.SetGuildCommands.Compile(nil, applicationID, guildID, commandCreates)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commands)
	return
}

func (s *ApplicationServiceImpl) UpdateGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID, commandUpdate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &command)
	return
}

func (s *ApplicationServiceImpl) DeleteGuildCommand(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) Error {
	compiledRoute, err := route.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(ctx, compiledRoute, nil, nil)
}

func (s *ApplicationServiceImpl) GetGuildCommandsPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake) (commandsPerms []discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commandsPerms)
	return
}

func (s *ApplicationServiceImpl) GetGuildCommandPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (commandPerms *discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commandPerms)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommandsPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) (commandsPerms []discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID, commandPermissions)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commandsPerms)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommandPermissions(ctx context.Context, applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions ...discord.CommandPermission) (commandPerms *discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID, commandPermissions)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &commandPerms)
	return
}
