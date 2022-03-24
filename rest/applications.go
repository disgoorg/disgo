package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Applications = (*applicationsImpl)(nil)

func NewApplications(restClient Client) Applications {
	return &applicationsImpl{restClient: restClient}
}

type Applications interface {
	GetGlobalCommands(applicationID snowflake.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGlobalCommand(applicationID snowflake.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGlobalCommands(applicationID snowflake.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) error

	GetGuildCommands(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, command discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGuildCommands(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commands []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, command discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) error

	GetGuildCommandsPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.ApplicationCommandPermissions, error)
	GetGuildCommandPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...RequestOpt) ([]discord.ApplicationCommandPermissions, error)
	SetGuildCommandPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)
}

type applicationsImpl struct {
	restClient Client
}

func (s *applicationsImpl) GetGlobalCommands(applicationID snowflake.Snowflake, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, nil, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGlobalCommand(applicationID snowflake.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandCreate, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGlobalCommands(applicationID snowflake.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGlobalCommand(applicationID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommands(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, nil, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandCreate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGuildCommands(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.restClient.Do(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGuildCommand(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommandsPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &commandsPerms, opts...)
	return
}

func (s *applicationsImpl) GetGuildCommandPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &commandPerms, opts...)
	return
}

func (s *applicationsImpl) SetGuildCommandsPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, commandPermissions, &commandsPerms, opts...)
	return
}

func (s *applicationsImpl) SetGuildCommandPermissions(applicationID snowflake.Snowflake, guildID snowflake.Snowflake, commandID snowflake.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, discord.ApplicationCommandPermissionsSet{Permissions: commandPermissions}, &commandPerms, opts...)
	return
}

func unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands []discord.UnmarshalApplicationCommand) []discord.ApplicationCommand {
	commands := make([]discord.ApplicationCommand, len(unmarshalCommands))
	for i := range unmarshalCommands {
		commands[i] = unmarshalCommands[i].ApplicationCommand
	}
	return commands
}
