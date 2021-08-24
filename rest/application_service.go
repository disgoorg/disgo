package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ ApplicationService = (*ApplicationServiceImpl)(nil)

func NewApplicationService(client Client) ApplicationService {
	return &ApplicationServiceImpl{restClient: client}
}

type ApplicationService interface {
	Service
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, Error)
	GetAuthorizationInfo(opts ...RequestOpt) (*discord.AuthorizationInformation, Error)

	GetGlobalCommands(applicationID discord.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, Error)
	GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	CreateGlobalCommand(applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	SetGlobalCommands(applicationID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, Error)
	UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) Error

	GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, Error)
	GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, command discord.ApplicationCommandCreate, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commands []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, Error)
	UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate, opts ...RequestOpt) (*discord.ApplicationCommand, Error)
	DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) Error

	GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) ([]discord.GuildCommandPermissions, Error)
	GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (*discord.GuildCommandPermissions, Error)
	SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions []discord.GuildCommandPermissionsSet, opts ...RequestOpt) ([]discord.GuildCommandPermissions, Error)
	SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions []discord.CommandPermission, opts ...RequestOpt) (*discord.GuildCommandPermissions, Error)
}

type ApplicationServiceImpl struct {
	restClient Client
}

func (s *ApplicationServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *ApplicationServiceImpl) GetBotApplicationInfo(opts ...RequestOpt) (application *discord.Application, rErr Error) {
	compiledRoute, err := route.GetBotApplicationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &application, opts...)
	return
}

func (s *ApplicationServiceImpl) GetAuthorizationInfo(opts ...RequestOpt) (info *discord.AuthorizationInformation, rErr Error) {
	compiledRoute, err := route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &info, opts...)
	return
}

func (s *ApplicationServiceImpl) GetGlobalCommands(applicationID discord.Snowflake, opts ...RequestOpt) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commands, opts...)
	return
}

func (s *ApplicationServiceImpl) GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) CreateGlobalCommand(applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.CreateGlobalCommand.Compile(nil, applicationID, commandCreate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) SetGlobalCommands(applicationID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.SetGlobalCommands.Compile(nil, applicationID, commandCreates)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commands, opts...)
	return
}

func (s *ApplicationServiceImpl) UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.UpdateGlobalCommand.Compile(nil, applicationID, commandID, commandUpdate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ApplicationServiceImpl) GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) (command []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.CreateGuildCommand.Compile(nil, applicationID, guildID, commandCreate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.SetGuildCommands.Compile(nil, applicationID, guildID, commandCreates)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commands, opts...)
	return
}

func (s *ApplicationServiceImpl) UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command *discord.ApplicationCommand, rErr Error) {
	compiledRoute, err := route.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID, commandUpdate)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &command, opts...)
	return
}

func (s *ApplicationServiceImpl) DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ApplicationServiceImpl) GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) (commandsPerms []discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commandsPerms, opts...)
	return
}

func (s *ApplicationServiceImpl) GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (commandPerms *discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commandPerms, opts...)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions []discord.GuildCommandPermissionsSet, opts ...RequestOpt) (commandsPerms []discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID, commandPermissions)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commandsPerms, opts...)
	return
}

func (s *ApplicationServiceImpl) SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions []discord.CommandPermission, opts ...RequestOpt) (commandPerms *discord.GuildCommandPermissions, rErr Error) {
	compiledRoute, err := route.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID, commandPermissions)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &commandPerms, opts...)
	return
}
