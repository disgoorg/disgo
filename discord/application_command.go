package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

type ApplicationCommandType int

const (
	ApplicationCommandTypeSlash = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
)

type ApplicationCommand interface {
	json.Marshaler
	ID() snowflake.Snowflake
	Type() ApplicationCommandType
	ApplicationID() snowflake.Snowflake
	GuildID() *snowflake.Snowflake
	Name() string
	NameLocalizations() map[Locale]string
	NameLocalized() string
	DefaultPermission() bool
	Version() snowflake.Snowflake
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

	default:
		err = fmt.Errorf("unkown application command with type %d received", cType.Type)
	}

	if err != nil {
		return err
	}

	u.ApplicationCommand = applicationCommand
	return nil
}

var _ ApplicationCommand = (*SlashCommand)(nil)

type SlashCommand struct {
	id                       snowflake.Snowflake
	applicationID            snowflake.Snowflake
	guildID                  *snowflake.Snowflake
	name                     string
	nameLocalizations        map[Locale]string
	nameLocalized            string
	Description              string
	DescriptionLocalizations map[Locale]string
	DescriptionLocalized     string
	Options                  []ApplicationCommandOption
	defaultPermission        bool
	version                  snowflake.Snowflake
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
	c.defaultPermission = v.DefaultPermission
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
		DefaultPermission:        c.defaultPermission,
		Version:                  c.version,
	})
}

func (c SlashCommand) ID() snowflake.Snowflake {
	return c.id
}

func (SlashCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (c SlashCommand) ApplicationID() snowflake.Snowflake {
	return c.applicationID
}

func (c SlashCommand) GuildID() *snowflake.Snowflake {
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

func (c SlashCommand) DefaultPermission() bool {
	return c.defaultPermission
}

func (c SlashCommand) Version() snowflake.Snowflake {
	return c.version
}

func (SlashCommand) applicationCommand() {}

var _ ApplicationCommand = (*UserCommand)(nil)

type UserCommand struct {
	id                snowflake.Snowflake
	applicationID     snowflake.Snowflake
	guildID           *snowflake.Snowflake
	name              string
	nameLocalizations map[Locale]string
	nameLocalized     string
	defaultPermission bool
	version           snowflake.Snowflake
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
	c.defaultPermission = v.DefaultPermission
	c.version = v.Version
	return nil
}

func (c UserCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawContextCommand{
		ID:                c.id,
		Type:              c.Type(),
		ApplicationID:     c.applicationID,
		GuildID:           c.guildID,
		Name:              c.name,
		NameLocalizations: c.nameLocalizations,
		NameLocalized:     c.nameLocalized,
		DefaultPermission: c.defaultPermission,
		Version:           c.version,
	})
}

func (c UserCommand) ID() snowflake.Snowflake {
	return c.id
}

func (c UserCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (c UserCommand) ApplicationID() snowflake.Snowflake {
	return c.applicationID
}

func (c UserCommand) GuildID() *snowflake.Snowflake {
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

func (c UserCommand) DefaultPermission() bool {
	return c.defaultPermission
}

func (c UserCommand) Version() snowflake.Snowflake {
	return c.version
}

func (UserCommand) applicationCommand() {}

var _ ApplicationCommand = (*MessageCommand)(nil)

type MessageCommand struct {
	id                snowflake.Snowflake
	applicationID     snowflake.Snowflake
	guildID           *snowflake.Snowflake
	name              string
	nameLocalizations map[Locale]string
	nameLocalized     string
	defaultPermission bool
	version           snowflake.Snowflake
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
	c.defaultPermission = v.DefaultPermission
	c.version = v.Version
	return nil
}

func (c MessageCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawContextCommand{
		ID:                c.id,
		Type:              c.Type(),
		ApplicationID:     c.applicationID,
		GuildID:           c.guildID,
		Name:              c.name,
		NameLocalizations: c.nameLocalizations,
		NameLocalized:     c.nameLocalized,
		DefaultPermission: c.defaultPermission,
		Version:           c.version,
	})
}

func (c MessageCommand) ID() snowflake.Snowflake {
	return c.id
}

func (MessageCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (c MessageCommand) ApplicationID() snowflake.Snowflake {
	return c.applicationID
}

func (c MessageCommand) GuildID() *snowflake.Snowflake {
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

func (c MessageCommand) DefaultPermission() bool {
	return c.defaultPermission
}

func (c MessageCommand) Version() snowflake.Snowflake {
	return c.version
}

func (MessageCommand) applicationCommand() {}
