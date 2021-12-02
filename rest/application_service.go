package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service            = (*applicationServiceImpl)(nil)
	_ ApplicationService = (*applicationServiceImpl)(nil)
)

func NewApplicationService(restClient Client) ApplicationService {
	return &applicationServiceImpl{restClient: restClient}
}

type ApplicationService interface {
	Service
	GetGlobalCommands(applicationID discord.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGlobalCommand(applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGlobalCommands(applicationID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) error

	GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, command discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commands []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) error

	GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommandPermissions, error)
	GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...RequestOpt) ([]discord.ApplicationCommandPermissions, error)
	SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)
}

type applicationServiceImpl struct {
	restClient Client
}

func (s *applicationServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *applicationServiceImpl) GetGlobalCommands(applicationID discord.Snowflake, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationServiceImpl) GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, nil, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) CreateGlobalCommand(applicationID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandCreate, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) SetGlobalCommands(applicationID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationServiceImpl) UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *applicationServiceImpl) GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationServiceImpl) GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, nil, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandCreate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationServiceImpl) UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.DoBot(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationServiceImpl) DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *applicationServiceImpl) GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &commandsPerms, opts...)
	return
}

func (s *applicationServiceImpl) GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &commandPerms, opts...)
	return
}

func (s *applicationServiceImpl) SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, commandPermissions, &commandsPerms, opts...)
	return
}

func (s *applicationServiceImpl) SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, commandPermissions, &commandPerms, opts...)
	return
}

func unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands []discord.UnmarshalApplicationCommand) []discord.ApplicationCommand {
	commands := make([]discord.ApplicationCommand, len(unmarshalCommands))
	for i := range unmarshalCommands {
		commands[i] = unmarshalCommands[i].ApplicationCommand
	}
	return commands
}
