package rest

import (
	"github.com/disgoorg/disgo/internal/slicehelper"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Applications = (*applicationsImpl)(nil)

func NewApplications(client Client) Applications {
	return &applicationsImpl{client: client}
}

type Applications interface {
	GetCurrentApplication(opts ...RequestOpt) (*discord.Application, error)
	UpdateCurrentApplication(applicationUpdate discord.ApplicationUpdate, opts ...RequestOpt) (*discord.Application, error)

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

	GetApplicationRoleConnectionMetadata(applicationID snowflake.ID, opts ...RequestOpt) ([]discord.ApplicationRoleConnectionMetadata, error)
	UpdateApplicationRoleConnectionMetadata(applicationID snowflake.ID, newRecords []discord.ApplicationRoleConnectionMetadata, opts ...RequestOpt) ([]discord.ApplicationRoleConnectionMetadata, error)

	GetEntitlements(applicationID snowflake.ID, userID snowflake.ID, guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, excludeEnded bool, skuIDs []snowflake.ID, opts ...RequestOpt) ([]discord.Entitlement, error)
	CreateTestEntitlement(applicationID snowflake.ID, entitlementCreate discord.TestEntitlementCreate, opts ...RequestOpt) (*discord.Entitlement, error)
	DeleteTestEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error
	ConsumeEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error

	GetSKUs(applicationID snowflake.ID, opts ...RequestOpt) ([]discord.SKU, error)
}

type applicationsImpl struct {
	client Client
}

func (s *applicationsImpl) GetCurrentApplication(opts ...RequestOpt) (application *discord.Application, err error) {
	err = s.client.Do(GetCurrentApplication.Compile(nil), nil, &application, opts...)
	return
}

func (s *applicationsImpl) UpdateCurrentApplication(applicationUpdate discord.ApplicationUpdate, opts ...RequestOpt) (application *discord.Application, err error) {
	err = s.client.Do(UpdateCurrentApplication.Compile(nil), applicationUpdate, &application, opts...)
	return
}

func (s *applicationsImpl) GetGlobalCommands(applicationID snowflake.ID, withLocalizations bool, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(GetGlobalCommands.Compile(discord.QueryValues{"with_localizations": withLocalizations}, applicationID), nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(GetGlobalCommand.Compile(nil, applicationID, commandID), nil, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGlobalCommand(applicationID snowflake.ID, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(CreateGlobalCommand.Compile(nil, applicationID), commandCreate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGlobalCommands(applicationID snowflake.ID, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(SetGlobalCommands.Compile(nil, applicationID), commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(UpdateGlobalCommand.Compile(nil, applicationID, commandID), commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGlobalCommand(applicationID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGlobalCommand.Compile(nil, applicationID, commandID), nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, withLocalizations bool, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(GetGuildCommands.Compile(discord.QueryValues{"with_localizations": withLocalizations}, applicationID, guildID), nil, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) GetGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(GetGuildCommand.Compile(nil, applicationID, guildID, commandID), nil, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) CreateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandCreate discord.ApplicationCommandCreate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(CreateGuildCommand.Compile(nil, applicationID, guildID), commandCreate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) SetGuildCommands(applicationID snowflake.ID, guildID snowflake.ID, commandCreates []discord.ApplicationCommandCreate, opts ...RequestOpt) (commands []discord.ApplicationCommand, err error) {
	var unmarshalCommands []discord.UnmarshalApplicationCommand
	err = s.client.Do(SetGuildCommands.Compile(nil, applicationID, guildID), commandCreates, &unmarshalCommands, opts...)
	if err == nil {
		commands = unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands)
	}
	return
}

func (s *applicationsImpl) UpdateGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandUpdate discord.ApplicationCommandUpdate, opts ...RequestOpt) (command discord.ApplicationCommand, err error) {
	var unmarshalCommand discord.UnmarshalApplicationCommand
	err = s.client.Do(UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID), commandUpdate, &unmarshalCommand, opts...)
	if err == nil {
		command = unmarshalCommand.ApplicationCommand
	}
	return
}

func (s *applicationsImpl) DeleteGuildCommand(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID), nil, nil, opts...)
}

func (s *applicationsImpl) GetGuildCommandsPermissions(applicationID snowflake.ID, guildID snowflake.ID, opts ...RequestOpt) (commandsPerms []discord.ApplicationCommandPermissions, err error) {
	err = s.client.Do(GetGuildCommandsPermissions.Compile(nil, applicationID, guildID), nil, &commandsPerms, opts...)
	return
}

func (s *applicationsImpl) GetGuildCommandPermissions(applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	err = s.client.Do(GetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID), nil, &commandPerms, opts...)
	return
}

func (s *applicationsImpl) GetApplicationRoleConnectionMetadata(applicationID snowflake.ID, opts ...RequestOpt) (metadata []discord.ApplicationRoleConnectionMetadata, err error) {
	err = s.client.Do(GetApplicationRoleConnectionMetadata.Compile(nil, applicationID), nil, &metadata, opts...)
	return
}

func (s *applicationsImpl) UpdateApplicationRoleConnectionMetadata(applicationID snowflake.ID, newRecords []discord.ApplicationRoleConnectionMetadata, opts ...RequestOpt) (metadata []discord.ApplicationRoleConnectionMetadata, err error) {
	err = s.client.Do(UpdateApplicationRoleConnectionMetadata.Compile(nil, applicationID), newRecords, &metadata, opts...)
	return
}

func (s *applicationsImpl) GetEntitlements(applicationID snowflake.ID, userID snowflake.ID, guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, excludeEnded bool, skuIDs []snowflake.ID, opts ...RequestOpt) (entitlements []discord.Entitlement, err error) {
	queryValues := discord.QueryValues{
		"exclude_ended": excludeEnded,
		"sku_ids":       slicehelper.JoinSnowflakes(skuIDs),
	}
	if userID != 0 {
		queryValues["user_id"] = userID
	}
	if guildID != 0 {
		queryValues["guild_id"] = guildID
	}
	if before != 0 {
		queryValues["before"] = before
	}
	if after != 0 {
		queryValues["after"] = after
	}
	if limit != 0 {
		queryValues["limit"] = limit
	}
	err = s.client.Do(GetEntitlements.Compile(queryValues, applicationID), nil, &entitlements, opts...)
	return
}

func (s *applicationsImpl) CreateTestEntitlement(applicationID snowflake.ID, entitlementCreate discord.TestEntitlementCreate, opts ...RequestOpt) (entitlement *discord.Entitlement, err error) {
	err = s.client.Do(CreateTestEntitlement.Compile(nil, applicationID), entitlementCreate, &entitlement, opts...)
	return
}

func (s *applicationsImpl) DeleteTestEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteTestEntitlement.Compile(nil, applicationID, entitlementID), nil, nil, opts...)
}

func (s *applicationsImpl) ConsumeEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(ConsumeEntitlement.Compile(nil, applicationID, entitlementID), nil, nil, opts...)
}

func (s *applicationsImpl) GetSKUs(applicationID snowflake.ID, opts ...RequestOpt) (skus []discord.SKU, err error) {
	err = s.client.Do(GetSKUs.Compile(nil, applicationID), nil, &skus, opts...)
	return
}

func unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands []discord.UnmarshalApplicationCommand) []discord.ApplicationCommand {
	commands := make([]discord.ApplicationCommand, len(unmarshalCommands))
	for i := range unmarshalCommands {
		commands[i] = unmarshalCommands[i].ApplicationCommand
	}
	return commands
}
