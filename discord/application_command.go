package discord

import (
	"fmt"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

type ApplicationCommandType int

const (
	ApplicationCommandTypeSlash ApplicationCommandType = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
	ApplicationCommandTypePrimaryEntryPoint
)

type ApplicationCommand interface {
	json.Marshaler
	ID() snowflake.ID
	Type() ApplicationCommandType
	ApplicationID() snowflake.ID
	GuildID() *snowflake.ID
	Name() string
	NameLocalizations() map[Locale]string
	NameLocalized() string
	DefaultMemberPermissions() Permissions
	DMPermission() bool
	Version() snowflake.ID
	CreatedAt() time.Time
	NSFW() bool
	IntegrationTypes() []ApplicationIntegrationType
	Contexts() []InteractionContextType
	applicationCommand()
}

type UnmarshalApplicationCommand struct {
	ApplicationCommand
}

func (u *UnmarshalApplicationCommand) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		applicationCommand ApplicationCommand
		err                error
	)

	switch cType.Type {
	case ApplicationCommandTypeSlash:
		var v SlashCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeUser:
		var v UserCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeMessage:
		var v MessageCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypePrimaryEntryPoint:
		var v EntryPointCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	default:
		err = fmt.Errorf("unknown application command with type %d received", cType.Type)
	}

	if err != nil {
		return err
	}

	u.ApplicationCommand = applicationCommand
	return nil
}

var _ ApplicationCommand = (*SlashCommand)(nil)

type SlashCommand struct {
	id                       snowflake.ID
	applicationID            snowflake.ID
	guildID                  *snowflake.ID
	name                     string
	nameLocalizations        map[Locale]string
	nameLocalized            string
	Description              string
	DescriptionLocalizations map[Locale]string
	DescriptionLocalized     string
	Options                  []ApplicationCommandOption
	defaultMemberPermissions Permissions
	dmPermission             bool
	nsfw                     bool
	integrationTypes         []ApplicationIntegrationType
	contexts                 []InteractionContextType
	version                  snowflake.ID
}

func (c *SlashCommand) UnmarshalJSON(data []byte) error {
	var v rawSlashCommand
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.applicationID = v.ApplicationID
	c.guildID = v.GuildID
	c.name = v.Name
	c.nameLocalizations = v.NameLocalizations
	c.nameLocalized = v.NameLocalized
	c.Description = v.Description
	c.DescriptionLocalizations = v.DescriptionLocalizations
	c.DescriptionLocalized = v.DescriptionLocalized
	c.Options = v.Options
	c.defaultMemberPermissions = v.DefaultMemberPermissions
	c.dmPermission = v.DMPermission
	c.nsfw = v.NSFW
	c.integrationTypes = v.IntegrationTypes
	c.contexts = v.Contexts
	c.version = v.Version
	return nil
}

func (c SlashCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawSlashCommand{
		ID:                       c.id,
		Type:                     c.Type(),
		ApplicationID:            c.applicationID,
		GuildID:                  c.guildID,
		Name:                     c.name,
		NameLocalizations:        c.nameLocalizations,
		NameLocalized:            c.nameLocalized,
		Description:              c.Description,
		DescriptionLocalizations: c.DescriptionLocalizations,
		DescriptionLocalized:     c.DescriptionLocalized,
		Options:                  c.Options,
		DefaultMemberPermissions: c.defaultMemberPermissions,
		DMPermission:             c.dmPermission,
		NSFW:                     c.nsfw,
		IntegrationTypes:         c.integrationTypes,
		Contexts:                 c.contexts,
		Version:                  c.version,
	})
}

func (c SlashCommand) ID() snowflake.ID {
	return c.id
}

func (SlashCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (c SlashCommand) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c SlashCommand) GuildID() *snowflake.ID {
	return c.guildID
}

func (c SlashCommand) Name() string {
	return c.name
}

func (c SlashCommand) NameLocalizations() map[Locale]string {
	return c.nameLocalizations
}

func (c SlashCommand) NameLocalized() string {
	return c.nameLocalized
}

func (c SlashCommand) DefaultMemberPermissions() Permissions {
	return c.defaultMemberPermissions
}

func (c SlashCommand) DMPermission() bool {
	return c.dmPermission
}

func (c SlashCommand) NSFW() bool {
	return c.nsfw
}

func (c SlashCommand) IntegrationTypes() []ApplicationIntegrationType {
	return c.integrationTypes
}

func (c SlashCommand) Contexts() []InteractionContextType {
	return c.contexts
}

func (c SlashCommand) Version() snowflake.ID {
	return c.version
}

func (c SlashCommand) CreatedAt() time.Time {
	return c.id.Time()
}

func (c SlashCommand) Mention() string {
	return SlashCommandMention(c.id, c.name)
}

func (SlashCommand) applicationCommand() {}

var _ ApplicationCommand = (*UserCommand)(nil)

type UserCommand struct {
	id                       snowflake.ID
	applicationID            snowflake.ID
	guildID                  *snowflake.ID
	name                     string
	nameLocalizations        map[Locale]string
	nameLocalized            string
	defaultMemberPermissions Permissions
	dmPermission             bool
	nsfw                     bool
	integrationTypes         []ApplicationIntegrationType
	contexts                 []InteractionContextType
	version                  snowflake.ID
}

func (c *UserCommand) UnmarshalJSON(data []byte) error {
	var v rawContextCommand
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.applicationID = v.ApplicationID
	c.guildID = v.GuildID
	c.name = v.Name
	c.nameLocalizations = v.NameLocalizations
	c.nameLocalized = v.NameLocalized
	c.defaultMemberPermissions = v.DefaultMemberPermissions
	c.dmPermission = v.DMPermission
	c.nsfw = v.NSFW
	c.integrationTypes = v.IntegrationTypes
	c.contexts = v.Contexts
	c.version = v.Version
	return nil
}

func (c UserCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawContextCommand{
		ID:                       c.id,
		Type:                     c.Type(),
		ApplicationID:            c.applicationID,
		GuildID:                  c.guildID,
		Name:                     c.name,
		NameLocalizations:        c.nameLocalizations,
		NameLocalized:            c.nameLocalized,
		DefaultMemberPermissions: c.defaultMemberPermissions,
		DMPermission:             c.dmPermission,
		NSFW:                     c.nsfw,
		IntegrationTypes:         c.integrationTypes,
		Contexts:                 c.contexts,
		Version:                  c.version,
	})
}

func (c UserCommand) ID() snowflake.ID {
	return c.id
}

func (c UserCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (c UserCommand) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c UserCommand) GuildID() *snowflake.ID {
	return c.guildID
}

func (c UserCommand) Name() string {
	return c.name
}

func (c UserCommand) NameLocalizations() map[Locale]string {
	return c.nameLocalizations
}

func (c UserCommand) NameLocalized() string {
	return c.nameLocalized
}

func (c UserCommand) DefaultMemberPermissions() Permissions {
	return c.defaultMemberPermissions
}

func (c UserCommand) DMPermission() bool {
	return c.dmPermission
}

func (c UserCommand) NSFW() bool {
	return c.nsfw
}

func (c UserCommand) IntegrationTypes() []ApplicationIntegrationType {
	return c.integrationTypes
}

func (c UserCommand) Contexts() []InteractionContextType {
	return c.contexts
}

func (c UserCommand) Version() snowflake.ID {
	return c.version
}

func (c UserCommand) CreatedAt() time.Time {
	return c.id.Time()
}

func (UserCommand) applicationCommand() {}

var _ ApplicationCommand = (*MessageCommand)(nil)

type MessageCommand struct {
	id                       snowflake.ID
	applicationID            snowflake.ID
	guildID                  *snowflake.ID
	name                     string
	nameLocalizations        map[Locale]string
	nameLocalized            string
	defaultMemberPermissions Permissions
	dmPermission             bool
	nsfw                     bool
	integrationTypes         []ApplicationIntegrationType
	contexts                 []InteractionContextType
	version                  snowflake.ID
}

func (c *MessageCommand) UnmarshalJSON(data []byte) error {
	var v rawContextCommand
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.applicationID = v.ApplicationID
	c.guildID = v.GuildID
	c.name = v.Name
	c.nameLocalizations = v.NameLocalizations
	c.nameLocalized = v.NameLocalized
	c.defaultMemberPermissions = v.DefaultMemberPermissions
	c.dmPermission = v.DMPermission
	c.nsfw = v.NSFW
	c.integrationTypes = v.IntegrationTypes
	c.contexts = v.Contexts
	c.version = v.Version
	return nil
}

func (c MessageCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawContextCommand{
		ID:                       c.id,
		Type:                     c.Type(),
		ApplicationID:            c.applicationID,
		GuildID:                  c.guildID,
		Name:                     c.name,
		NameLocalizations:        c.nameLocalizations,
		NameLocalized:            c.nameLocalized,
		DefaultMemberPermissions: c.defaultMemberPermissions,
		DMPermission:             c.dmPermission,
		NSFW:                     c.nsfw,
		IntegrationTypes:         c.integrationTypes,
		Contexts:                 c.contexts,
		Version:                  c.version,
	})
}

func (c MessageCommand) ID() snowflake.ID {
	return c.id
}

func (MessageCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (c MessageCommand) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c MessageCommand) GuildID() *snowflake.ID {
	return c.guildID
}

func (c MessageCommand) Name() string {
	return c.name
}

func (c MessageCommand) NameLocalizations() map[Locale]string {
	return c.nameLocalizations
}

func (c MessageCommand) NameLocalized() string {
	return c.nameLocalized
}

func (c MessageCommand) DefaultMemberPermissions() Permissions {
	return c.defaultMemberPermissions
}

func (c MessageCommand) DMPermission() bool {
	return c.dmPermission
}

func (c MessageCommand) NSFW() bool {
	return c.nsfw
}

func (c MessageCommand) IntegrationTypes() []ApplicationIntegrationType {
	return c.integrationTypes
}

func (c MessageCommand) Contexts() []InteractionContextType {
	return c.contexts
}

func (c MessageCommand) Version() snowflake.ID {
	return c.version
}

func (c MessageCommand) CreatedAt() time.Time {
	return c.id.Time()
}

func (MessageCommand) applicationCommand() {}

var _ ApplicationCommand = (*EntryPointCommand)(nil)

type EntryPointCommand struct {
	id                       snowflake.ID
	applicationID            snowflake.ID
	guildID                  *snowflake.ID
	name                     string
	nameLocalizations        map[Locale]string
	nameLocalized            string
	defaultMemberPermissions Permissions
	dmPermission             bool
	nsfw                     bool
	integrationTypes         []ApplicationIntegrationType
	contexts                 []InteractionContextType
	version                  snowflake.ID
	Handler                  EntryPointCommandHandlerType
}

func (c *EntryPointCommand) UnmarshalJSON(data []byte) error {
	var v rawEntryPointCommand
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.applicationID = v.ApplicationID
	c.guildID = v.GuildID
	c.name = v.Name
	c.nameLocalizations = v.NameLocalizations
	c.nameLocalized = v.NameLocalized
	c.defaultMemberPermissions = v.DefaultMemberPermissions
	c.dmPermission = v.DMPermission
	c.nsfw = v.NSFW
	c.integrationTypes = v.IntegrationTypes
	c.contexts = v.Contexts
	c.version = v.Version
	c.Handler = v.Handler
	return nil
}

func (c EntryPointCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawEntryPointCommand{
		ID:                       c.id,
		Type:                     c.Type(),
		ApplicationID:            c.applicationID,
		GuildID:                  c.guildID,
		Name:                     c.name,
		NameLocalizations:        c.nameLocalizations,
		NameLocalized:            c.nameLocalized,
		DefaultMemberPermissions: c.defaultMemberPermissions,
		DMPermission:             c.dmPermission,
		NSFW:                     c.nsfw,
		IntegrationTypes:         c.integrationTypes,
		Contexts:                 c.contexts,
		Version:                  c.version,
		Handler:                  c.Handler,
	})
}

func (c EntryPointCommand) ID() snowflake.ID {
	return c.id
}

func (EntryPointCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypePrimaryEntryPoint
}

func (c EntryPointCommand) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c EntryPointCommand) GuildID() *snowflake.ID {
	return c.guildID
}

func (c EntryPointCommand) Name() string {
	return c.name
}

func (c EntryPointCommand) NameLocalizations() map[Locale]string {
	return c.nameLocalizations
}

func (c EntryPointCommand) NameLocalized() string {
	return c.nameLocalized
}

func (c EntryPointCommand) DefaultMemberPermissions() Permissions {
	return c.defaultMemberPermissions
}

func (c EntryPointCommand) DMPermission() bool {
	return c.dmPermission
}

func (c EntryPointCommand) NSFW() bool {
	return c.nsfw
}

func (c EntryPointCommand) IntegrationTypes() []ApplicationIntegrationType {
	return c.integrationTypes
}

func (c EntryPointCommand) Contexts() []InteractionContextType {
	return c.contexts
}

func (c EntryPointCommand) Version() snowflake.ID {
	return c.version
}

func (c EntryPointCommand) CreatedAt() time.Time {
	return c.id.Time()
}

func (EntryPointCommand) applicationCommand() {}

type EntryPointCommandHandlerType int

const (
	EntryPointCommandHandlerTypeAppHandler EntryPointCommandHandlerType = iota + 1
	EntryPointCommandHandlerTypeDiscordLaunchActivity
)
