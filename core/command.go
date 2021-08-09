package core

import "github.com/DisgoOrg/disgo/discord"

type Command struct {
	discord.Command
	Disgo            Disgo
	GuildPermissions map[discord.Snowflake]*GuildCommandPermissions
}

type GuildCommandPermissions struct {
	discord.GuildCommandPermissions
	Disgo Disgo
}

// Guild returns the Guild the Command is from the Cache or nil if it is a global Command
func (c Command) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Disgo.Cache().GuildCache().Get(*c.GuildID)
}

// FromGuild returns true if this is a guild Command else false
func (c Command) FromGuild() bool {
	return c.GuildID == nil
}

// ToCreate return the CommandCreate for this Command
func (c Command) ToCreate() discord.CommandCreate {
	return discord.CommandCreate{
		Name:              c.Name,
		Description:       c.Description,
		DefaultPermission: c.DefaultPermission,
		Options:           c.Options,
	}
}

// Update updates the current Command with the given fields
func (c Command) Update(command discord.CommandUpdate) error {
	var rC *Command
	var err error
	if c.GuildID == nil {
		rC, err = c.Disgo.RestServices().UpdateGlobalCommand(c.Disgo.ApplicationID(), c.ID, command)

	} else {
		rC, err = c.Disgo.RestServices().UpdateGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID, command)
	}
	if err != nil {
		return err
	}
	*c = *rC
	return nil
}

// SetPermissions sets the GuildCommandPermissions for a specific Guild. this overrides all existing CommandPermission(s). thx discord for that
func (c Command) SetPermissions(guildID discord.Snowflake, permissions ...CommandPermission) error {
	_, err := c.Disgo.RestServices().SetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID, SetGuildCommandPermissions{Permissions: permissions})
	if err != nil {
		return err
	}
	return nil
}

// GetPermissions returns the GuildCommandPermissions for the specific Guild from the Cache
func (c Command) GetPermissions(guildID discord.Snowflake) *GuildCommandPermissions {
	return c.GuildPermissions[guildID]
}

// FetchPermissions fetched the GuildCommandPermissions for a specific Guild from discord
func (c *Command) FetchPermissions(guildID discord.Snowflake) (*GuildCommandPermissions, error) {
	perms, err := c.Disgo.RestServices().GetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID)
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
		return c.Disgo.RestServices().DeleteGlobalCommand(c.Disgo.ApplicationID(), c.ID)

	}
	return c.Disgo.RestServices().DeleteGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID)
}


