package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Applications = (*applicationsImpl)(nil)

func NewApplications(client Client) Applications {
	return &applicationsImpl{client: client}
}

type Applications interface {
	GetGlobalCommands(applicationID snowflake.ID, withLocalizations bool, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGlobalCommand(applicationID snowflake.ID, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGlobalCommands(applicationID snowflake.ID, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error

	GetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, withLocalizations bool, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	GetGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (discord.ApplicationCommand, error)
	CreateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, command discord.ApplicationCommandCreate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	SetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, commands []discord.ApplicationCommandCreate, opts ...RequestOpt) ([]discord.ApplicationCommand, error)
	UpdateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, command discord.ApplicationCommandUpdate, opts ...RequestOpt) (discord.ApplicationCommand, error)
	DeleteGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error

	GetGuildCommandsPermissions(applicationID snowflake.ID, guildID snowflake.ID, opts ...RequestOpt) ([]discord.ApplicationCommandPermissions, error)
	GetGuildCommandPermissions(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)
}

type applicationsImpl struct {
	client Client
}

func (s *applicationsImpl) GetGlobalCommands(applicationID snowflake.ID, withLocalizations bool, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommands.Compile(route.QueryValues{"with_localizations": withLocalizations}, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, nil, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGlobalCommand(applicationID snowflake.ID, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandCreate, &command, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGlobalCommands(applicationID snowflake.ID, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, withLocalizations bool, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommands.Compile(route.QueryValues{"with_localizations": withLocalizations}, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, nil, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandCreate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(compiledRoute, commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommandsPermissions(applicationID snowflake.ID, guildID snowflake.ID, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &commandsPerms, opts...)
	return
}

func (s *applicationsImpl) GetGuildCommandPermissions(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &commandPerms, opts...)
	return
}

func unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands []discord.UnmarshalApplicationCommand) []discord.ApplicationCommand {
	commands := make([]discord.ApplicationCommand, len(unmarshalCommands))
	for i := range unmarshalCommands {
		commands[i] = unmarshalCommands[i].ApplicationCommand
	}
	return commands
}
