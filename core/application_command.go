package core

import "github.com/DisgoOrg/disgo/discord"

type ApplicationCommand struct {
	discord.ApplicationCommand
	Disgo            Disgo
	// TODO: should we cache command perms per guild? extra cache & cache flag?
	//GuildPermissions map[discord.Snowflake]*GuildCommandPermissions
}

type GuildCommandPermissions struct {
	discord.GuildCommandPermissions
	Disgo Disgo
}

// Guild returns the Guild the ApplicationCommand is from the Cache or nil if it is a global ApplicationCommand
func (c *ApplicationCommand) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Disgo.Cache().GuildCache().Get(*c.GuildID)
}

// IsGlobal returns true if this is a global ApplicationCommand and false for a guild ApplicationCommand
func (c *ApplicationCommand) IsGlobal() bool {
	return c.GuildID == nil
}

// ToCreate return the ApplicationCommandCreate for this ApplicationCommand
func (c *ApplicationCommand) ToCreate() discord.ApplicationCommandCreate {
	return discord.ApplicationCommandCreate{
		Type:              c.Type,
		Name:              c.Name,
		Description:       c.Description,
		DefaultPermission: c.DefaultPermission,
		Options:           c.Options,
	}
}

// Update updates the current ApplicationCommand with the given fields
func (c *ApplicationCommand) Update(commandUpdate discord.ApplicationCommandUpdate) (*ApplicationCommand, error) {
	var command *discord.ApplicationCommand
	var err error
	if c.GuildID == nil {
		command, err = c.Disgo.RestServices().ApplicationService().UpdateGlobalCommand(c.Disgo.ApplicationID(), c.ID, commandUpdate)

	} else {
		command, err = c.Disgo.RestServices().ApplicationService().UpdateGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID, commandUpdate)
	}
	if err != nil {
		return nil, err
	}
	return c.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// SetPermissions sets the GuildCommandPermissions for a specific Guild. this overrides all existing CommandPermission(s). thx discord for that
func (c *ApplicationCommand) SetPermissions(guildID discord.Snowflake, commandPermissions ...discord.CommandPermission) (*GuildCommandPermissions, error) {
	permissions, err := c.Disgo.RestServices().ApplicationService().SetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID, commandPermissions...)
	if err != nil {
		return nil, err
	}
	return c.Disgo.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// GetPermissions fetched the GuildCommandPermissions for a specific Guild from discord
func (c *ApplicationCommand) GetPermissions(guildID discord.Snowflake) (*GuildCommandPermissions, error) {
	permissions, err := c.Disgo.RestServices().ApplicationService().GetGuildCommandPermissions(c.Disgo.ApplicationID(), guildID, c.ID)
	if err != nil {
		return nil, err
	}
	return c.Disgo.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// Delete deletes the ApplicationCommand from discord
func (c *ApplicationCommand) Delete() error {
	if c.GuildID == nil {
		return c.Disgo.RestServices().ApplicationService().DeleteGlobalCommand(c.Disgo.ApplicationID(), c.ID)
	}
	return c.Disgo.RestServices().ApplicationService().DeleteGuildCommand(c.Disgo.ApplicationID(), *c.GuildID, c.ID)
}
