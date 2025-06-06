package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/slicehelper"
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

	GetEntitlements(applicationID snowflake.ID, params GetEntitlementsParams, opts ...RequestOpt) ([]discord.Entitlement, error)
	GetEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) (*discord.Entitlement, error)
	CreateTestEntitlement(applicationID snowflake.ID, entitlementCreate discord.TestEntitlementCreate, opts ...RequestOpt) (*discord.Entitlement, error)
	DeleteTestEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error
	ConsumeEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) error

	GetApplicationEmojis(applicationID snowflake.ID, opts ...RequestOpt) ([]discord.Emoji, error)
	GetApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (*discord.Emoji, error)
	CreateApplicationEmoji(applicationID snowflake.ID, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (*discord.Emoji, error)
	UpdateApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (*discord.Emoji, error)
	DeleteApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) error

	GetActivityInstance(applicationID snowflake.ID, instanceID string, opts ...RequestOpt) (*discord.ActivityInstance, error)
}

// GetEntitlementsParams holds query parameters for Applications.GetEntitlements (https://discord.com/developers/docs/resources/entitlement#list-entitlements)
type GetEntitlementsParams struct {
	UserID         snowflake.ID
	SkuIDs         []snowflake.ID
	Before         int
	After          int
	Limit          int
	GuildID        snowflake.ID
	ExcludeEnded   bool
	ExcludeDeleted bool
}

func (p GetEntitlementsParams) ToQueryValues() discord.QueryValues {
	queryValues := discord.QueryValues{
		"exclude_ended":   p.ExcludeEnded,
		"exclude_deleted": p.ExcludeDeleted,
		"sku_ids":         slicehelper.JoinSnowflakes(p.SkuIDs),
	}
	if p.UserID != 0 {
		queryValues["user_id"] = p.UserID
	}
	if p.Before != 0 {
		queryValues["before"] = p.Before
	}
	if p.After != 0 {
		queryValues["after"] = p.After
	}
	if p.Limit != 0 {
		queryValues["limit"] = p.Limit
	}
	if p.GuildID != 0 {
		queryValues["guild_id"] = p.GuildID
	}
	return queryValues
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

func (s *applicationsImpl) GetEntitlements(applicationID snowflake.ID, params GetEntitlementsParams, opts ...RequestOpt) (entitlements []discord.Entitlement, err error) {
	err = s.client.Do(GetEntitlements.Compile(params.ToQueryValues(), applicationID), nil, &entitlements, opts...)
	return
}

func (s *applicationsImpl) GetEntitlement(applicationID snowflake.ID, entitlementID snowflake.ID, opts ...RequestOpt) (entitlement *discord.Entitlement, err error) {
	err = s.client.Do(GetEntitlement.Compile(nil, applicationID, entitlementID), nil, &entitlement, opts...)
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

func (s *applicationsImpl) GetApplicationEmojis(applicationID snowflake.ID, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var rs emojisResponse
	err = s.client.Do(GetApplicationEmojis.Compile(nil, applicationID), nil, &rs, opts...)
	if err == nil {
		emojis = rs.Items
	}
	return
}

func (s *applicationsImpl) GetApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	err = s.client.Do(GetApplicationEmoji.Compile(nil, applicationID, emojiID), nil, &emoji, opts...)
	return
}

func (s *applicationsImpl) CreateApplicationEmoji(applicationID snowflake.ID, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	err = s.client.Do(CreateApplicationEmoji.Compile(nil, applicationID), emojiCreate, &emoji, opts...)
	return
}

func (s *applicationsImpl) UpdateApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	err = s.client.Do(UpdateApplicationEmoji.Compile(nil, applicationID, emojiID), emojiUpdate, &emoji, opts...)
	return
}

func (s *applicationsImpl) DeleteApplicationEmoji(applicationID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteApplicationEmoji.Compile(nil, applicationID, emojiID), nil, nil, opts...)
}

func (s *applicationsImpl) GetActivityInstance(applicationID snowflake.ID, instanceID string, opts ...RequestOpt) (instance *discord.ActivityInstance, err error) {
	err = s.client.Do(GetActivityInstance.Compile(nil, applicationID, instanceID), nil, &instance, opts...)
	return
}

func unmarshalApplicationCommandsToApplicationCommands(unmarshalCommands []discord.UnmarshalApplicationCommand) []discord.ApplicationCommand {
	commands := make([]discord.ApplicationCommand, len(unmarshalCommands))
	for i := range unmarshalCommands {
		commands[i] = unmarshalCommands[i].ApplicationCommand
	}
	return commands
}

type emojisResponse struct {
	Items []discord.Emoji `json:"items"`
}
