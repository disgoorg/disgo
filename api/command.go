package api

import "errors"

var errNoDisgoInstance = errors.New("no disgo instance injected")

// Command is the base "command" model that belongs to an application.
type Command struct {
	Disgo             Disgo
	GuildPermissions  map[Snowflake]*GuildCommandPermissions
	GuildID           *Snowflake       `json:"guild_id"`
	ID                Snowflake        `json:"id,omitempty"`
	ApplicationID     Snowflake        `json:"application_id,omitempty"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	DefaultPermission bool             `json:"default_permission"`
	Options           []*CommandOption `json:"options,omitempty"`
}

// Guild returns the Guild the Command is from from the Cache or nil if it is a global Command
func (c Command) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Disgo.Cache().Guild(*c.GuildID)
}

// FromGuild returns true if this is a guild Command else false
func (c Command) FromGuild() bool {
	return c.GuildID == nil
}

// Create creates the Command either as global or guild Command depending on if GuildID is set
func (c *Command) Create() error {
	if c.Disgo == nil {
		return errNoDisgoInstance
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

// Fetch updates/fetches the current Command from discord
func (c *Command) Fetch() error {
	if c.Disgo == nil {
		return errNoDisgoInstance
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

// Update updates the current Command with the given fields
func (c *Command) Update(command UpdateCommand) error {
	if c.Disgo == nil {
		return errNoDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().EditGlobalCommand(c.Disgo.SelfUserID(), c.ID, command)

	} else {
		rC, err = c.Disgo.RestClient().EditGuildCommand(c.Disgo.SelfUserID(), *c.GuildID, c.ID, command)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

// SetPermissions sets the GuildCommandPermissions for a specific Guild. this overrides all existing CommandPermission(s). thx discord for that
func (c Command) SetPermissions(guildID Snowflake, permissions ...CommandPermission) (*GuildCommandPermissions, error) {
	perms, err := c.Disgo.RestClient().SetGuildCommandPermissions(c.Disgo.SelfUserID(), guildID, c.ID, SetGuildCommandPermissions{Permissions: permissions})
	if err != nil {
		return nil, err
	}
	c.GuildPermissions[guildID] = perms
	return perms, nil
}

// GetPermissions returns the GuildCommandPermissions for the specific Guild from the Cache
func (c Command) GetPermissions(guildID Snowflake) *GuildCommandPermissions {
	return c.GuildPermissions[guildID]
}

// FetchPermissions fetched the GuildCommandPermissions for a specific Guild from discord
func (c Command) FetchPermissions(guildID Snowflake) (*GuildCommandPermissions, error) {
	perms, err := c.Disgo.RestClient().GetGuildCommandPermissions(c.Disgo.SelfUserID(), guildID, c.ID)
	if err != nil {
		return nil, err
	}
	c.GuildPermissions[guildID] = perms
	return perms, nil
}

// Delete deletes the Command from discord
func (c Command) Delete() error {
	if c.Disgo == nil {
		return errNoDisgoInstance
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

// UpdateCommand is used to update an existing Command. all fields are optional
type UpdateCommand struct {
	Name              *string          `json:"name,omitempty"`
	Description       *string          `json:"description,omitempty"`
	DefaultPermission *bool            `json:"default_permission,omitempty"`
	Options           []*CommandOption `json:"options,omitempty"`
}
