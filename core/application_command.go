package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ApplicationCommand struct {
	discord.ApplicationCommand
	Bot     *Bot
	Options []ApplicationCommandOption
}

// Guild returns the Guild this ApplicationCommand is from or nil if this is a global ApplicationCommand.
// This will only check cached guilds!
func (c *ApplicationCommand) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Bot.Caches.GuildCache().Get(*c.GuildID)
}

// IsGlobal returns whether this ApplicationCommand is global
func (c *ApplicationCommand) IsGlobal() bool {
	return c.GuildID == nil
}

// ToCreate returns the discord.ApplicationCommandCreate for this ApplicationCommand, useful for cloning
func (c *ApplicationCommand) ToCreate() discord.ApplicationCommandCreate {
	return discord.ApplicationCommandCreate{
		Type:              c.Type,
		Name:              c.Name,
		Description:       c.Description,
		DefaultPermission: c.DefaultPermission,
		Options:           c.ApplicationCommand.Options,
	}
}

// Update updates this ApplicationCommand with the content provided in discord.ApplicationCommandUpdate
func (c *ApplicationCommand) Update(commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	var command *discord.ApplicationCommand
	var err rest.Error
	if c.GuildID == nil {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGlobalCommand(c.Bot.ApplicationID, c.ID, commandUpdate, opts...)

	} else {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, commandUpdate, opts...)
	}
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommand(*command), nil
}

// SetPermissions sets the ApplicationCommandPermissions of this ApplicationCommand for the specified Guild.
// This will override all existing ApplicationCommandPermissions!
func (c *ApplicationCommand) SetPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, rest.Error) {
	permissions, err := c.Bot.RestServices.ApplicationService().SetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// GetPermissions returns the ApplicationCommandPermissions of this ApplicationCommand for the specified Guild
func (c *ApplicationCommand) GetPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, rest.Error) {
	permissions, err := c.Bot.RestServices.ApplicationService().GetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// Delete deletes this ApplicationCommand
func (c *ApplicationCommand) Delete(opts ...rest.RequestOpt) rest.Error {
	if c.GuildID == nil {
		return c.Bot.RestServices.ApplicationService().DeleteGlobalCommand(c.Bot.ApplicationID, c.ID)
	}
	return c.Bot.RestServices.ApplicationService().DeleteGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, opts...)
}
