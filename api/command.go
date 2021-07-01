package api

import "errors"

var errNoDisgoInstance = errors.New("no disgo instance injected")

// Command is the base "command" model that belongs to an application.
type Command struct {
	Disgo             Disgo
	GuildPermissions  map[Snowflake]*GuildCommandPermissions
	GuildID           *Snowflake      `json:"guild_id"`
	ID                Snowflake       `json:"id,omitempty"`
	ApplicationID     Snowflake       `json:"application_id,omitempty"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	DefaultPermission bool            `json:"default_permission,omitempty"`
	Options           []CommandOption `json:"options,omitempty"`
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

// ToCreate return the CommandCreate for this Command
func (c *Command) ToCreate() CommandCreate {
	return CommandCreate{
		Name:              c.Name,
		Description:       c.Description,
		DefaultPermission: c.DefaultPermission,
		Options:           c.Options,
	}
}

// Fetch updates/fetches the current Command from discord
func (c *Command) Fetch() error {
	if c.Disgo == nil {
		return errNoDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().GetGlobalCommand(c.Disgo.ApplicationID(), c.ID)

	} else {
		rC, err = c.Disgo.RestClient().GetGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

// Update updates the current Command with the given fields
func (c *Command) Update(command CommandUpdate) error {
	if c.Disgo == nil {
		return errNoDisgoInstance
	}
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestClient().UpdateGlobalCommand(c.Disgo.ApplicationID(), c.ID, command)

	} else {
		rC, err = c.Disgo.RestClient().UpdateGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID, command)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

// SetPermissions sets the GuildCommandPermissions for a specific Guild. this overrides all existing CommandPermission(s). thx discord for that
func (c *Command) SetPermissions(guildID Snowflake, permissions ...CommandPermission) error {
	_, err := c.Disgo.RestClient().SetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID, SetGuildCommandPermissions{Permissions: permissions})
	if err != nil {
		return err
	}
	return nil
}

// GetPermissions returns the GuildCommandPermissions for the specific Guild from the Cache
func (c Command) GetPermissions(guildID Snowflake) *GuildCommandPermissions {
	return c.GuildPermissions[guildID]
}

// FetchPermissions fetched the GuildCommandPermissions for a specific Guild from discord
func (c *Command) FetchPermissions(guildID Snowflake) (*GuildCommandPermissions, error) {
	perms, err := c.Disgo.RestClient().GetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// Delete deletes the Command from discord
func (c Command) Delete() error {
	if c.Disgo == nil {
		return errNoDisgoInstance
	}
	if c.GuildID == nil {
		return c.Disgo.RestClient().DeleteGlobalCommand(c.Disgo.ApplicationID(), c.ID)

	}
	return c.Disgo.RestClient().DeleteGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID)
}

// CommandCreate is used to create an Command. all fields are optional
type CommandCreate struct {
	Name              string          `json:"name,omitempty"`
	Description       string          `json:"description,omitempty"`
	DefaultPermission bool            `json:"default_permission,omitempty"`
	Options           []CommandOption `json:"options,omitempty"`
}

// CommandUpdate is used to update an existing Command. all fields are optional
type CommandUpdate struct {
	Name              *string         `json:"name,omitempty"`
	Description       *string         `json:"description,omitempty"`
	DefaultPermission *bool           `json:"default_permission,omitempty"`
	Options           []CommandOption `json:"options,omitempty"`
}
