package api

import "errors"

var noDisgoInstance = errors.New("no disgo instance injected")

// Command is the base "command" model that belongs to an application.
type Command struct {
	Disgo             Disgo
	GuildID           *Snowflake
	ID                Snowflake        `json:"id,omitempty"`
	ApplicationID     Snowflake        `json:"application_id,omitempty"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	DefaultPermission bool             `json:"default_permission"`
	Options           []*CommandOption `json:"options,omitempty"`
}

func (c *Command) Create() error {
	if c.Disgo == nil {
		return noDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().CreateGlobalCommand(c.Disgo.SelfUserID(), *c)

	} else {
		rC, err = c.Disgo.RestClient().CreateGuildCommand(c.Disgo.SelfUserID(), *c.GuildID, *c)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

func (c *Command) Get() error {
	if c.Disgo == nil {
		return noDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().GetGlobalCommand(c.Disgo.SelfUserID(), c.ID)

	} else {
		rC, err = c.Disgo.RestClient().GetGuildCommand(c.Disgo.SelfUserID(), *c.GuildID, c.ID)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

func (c *Command) Update() error {
	if c.Disgo == nil {
		return noDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().EditGlobalCommand(c.Disgo.SelfUserID(), c.ID, *c)

	} else {
		rC, err = c.Disgo.RestClient().EditGuildCommand(c.Disgo.SelfUserID(), *c.GuildID, c.ID, *c)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

func (c Command) Delete() error {
	if c.Disgo == nil {
		return noDisgoInstance
	}
	if c.GuildID == nil {
		return c.Disgo.RestClient().DeleteGlobalCommand(c.Disgo.SelfUserID(), c.ID)

	}
	return c.Disgo.RestClient().DeleteGuildCommand(c.Disgo.SelfUserID(), *c.GuildID, c.ID)
}

// CommandOptionType specifies the type of the arguments used in Command.Options
type CommandOptionType int

// Constants for each slash command option type
const (
	CommandOptionTypeSubCommand CommandOptionType = iota + 1
	CommandOptionTypeSubCommandGroup
	CommandOptionTypeString
	CommandOptionTypeInteger
	CommandOptionTypeBoolean
	CommandOptionTypeUser
	CommandOptionTypeChannel
	CommandOptionTypeRole
)

// CommandOption are the arguments used in Command.Options
type CommandOption struct {
	Type        CommandOptionType `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required,omitempty"`
	Choices     []OptionChoice    `json:"choices,omitempty"`
	Options     []CommandOption   `json:"options,omitempty"`
}

// OptionChoice contains the data for a user using your command
type OptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// GuildCommandPermissions holds all permissions for a Command
type GuildCommandPermissions struct {
	ID            Snowflake           `json:"id"`
	ApplicationID Snowflake           `json:"application_id"`
	GuildID       Snowflake           `json:"guild_id"`
	Permissions   []CommandPermission `json:"permissions"`
}

// CommandPermissionType is the type of the CommandPermission
type CommandPermissionType int

const (
	CommandPermissionTypeRole = iota + 1
	CommandPermissionTypeUser
)

// CommandPermission holds a User or Role and if they are allowed to use the Command
type CommandPermission struct {
	ID         Snowflake             `json:"id"`
	Type       CommandPermissionType `json:"type"`
	Permission bool                  `json:"permission"`
}

// SetGuildCommandsPermissions holds a slice of SetGuildCommandPermissions
type SetGuildCommandsPermissions []SetGuildCommandPermissions

// SetGuildCommandPermissions is used to update CommandPermission ID should be omitted fro bulk update
type SetGuildCommandPermissions struct {
	ID          Snowflake           `json:"id,omitempty"`
	Permissions []CommandPermission `json:"permissions"`
}
