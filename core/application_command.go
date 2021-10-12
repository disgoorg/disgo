package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ApplicationCommand interface {
	discord.ApplicationCommand
}

type SlashCommand struct {
	discord.SlashCommand
	Bot *Bot
}

// Guild returns the Guild the ApplicationCommand is from the Caches or nil if it is a global ApplicationCommand
func (c *SlashCommand) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Bot.Caches.GuildCache().Get(*c.GuildID)
}

// IsGlobal returns true if this is a global ApplicationCommand and false for a guild ApplicationCommand
func (c *SlashCommand) IsGlobal() bool {

	return c.GuildID == nil
}

// ToCreate return the ApplicationCommandCreate for this ApplicationCommand
func (c *SlashCommand) ToCreate() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:              c.Name,
		Description:       c.Description,
		Options:           c.Options,
		DefaultPermission: c.DefaultPermission,
	}
}

// Update updates the current ApplicationCommand with the given fields
func (c *SlashCommand) Update(commandUpdate discord.SlashCommandUpdate, opts ...rest.RequestOpt) (*SlashCommand, error) {
	var command *discord.ApplicationCommand
	var err error
	if c.GuildID == nil {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGlobalCommand(c.Bot.ApplicationID, c.ID, commandUpdate, opts...)

	} else {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, commandUpdate, opts...)
	}
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateSlashCommand((*command).(discord.SlashCommand)), nil
}

// SetPermissions sets the ApplicationCommandPermissions for a specific Guild. this overrides all existing ApplicationCommandPermission(s). thx discord for that
func (c *SlashCommand) SetPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().SetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// GetPermissions fetched the ApplicationCommandPermissions for a specific Guild from discord
func (c *SlashCommand) GetPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().GetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// Delete deletes the ApplicationCommand from discord
func (c *SlashCommand) Delete(opts ...rest.RequestOpt) error {
	if c.GuildID == nil {
		return c.Bot.RestServices.ApplicationService().DeleteGlobalCommand(c.Bot.ApplicationID, c.ID)
	}
	return c.Bot.RestServices.ApplicationService().DeleteGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, opts...)
}

type UserCommand struct {
	discord.UserCommand
	Bot *Bot
}

// Guild returns the Guild the ApplicationCommand is from the Caches or nil if it is a global ApplicationCommand
func (c *UserCommand) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Bot.Caches.GuildCache().Get(*c.GuildID)
}

// IsGlobal returns true if this is a global ApplicationCommand and false for a guild ApplicationCommand
func (c *UserCommand) IsGlobal() bool {
	return c.GuildID == nil
}

// ToCreate return the ApplicationCommandCreate for this ApplicationCommand
func (c *UserCommand) ToCreate() discord.UserCommandCreate {
	return discord.UserCommandCreate{
		Name:              c.Name,
		DefaultPermission: c.DefaultPermission,
	}
}

// Update updates the current ApplicationCommand with the given fields
func (c *UserCommand) Update(commandUpdate discord.UserCommandUpdate, opts ...rest.RequestOpt) (*UserCommand, error) {
	var command *discord.ApplicationCommand
	var err error
	if c.GuildID == nil {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGlobalCommand(c.Bot.ApplicationID, c.ID, commandUpdate, opts...)

	} else {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, commandUpdate, opts...)
	}
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateUserCommand((*command).(discord.UserCommand)), nil
}

// SetPermissions sets the ApplicationCommandPermissions for a specific Guild. this overrides all existing ApplicationCommandPermission(s). thx discord for that
func (c *UserCommand) SetPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().SetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// GetPermissions fetched the ApplicationCommandPermissions for a specific Guild from discord
func (c *UserCommand) GetPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().GetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// Delete deletes the ApplicationCommand from discord
func (c *UserCommand) Delete(opts ...rest.RequestOpt) error {
	if c.GuildID == nil {
		return c.Bot.RestServices.ApplicationService().DeleteGlobalCommand(c.Bot.ApplicationID, c.ID)
	}
	return c.Bot.RestServices.ApplicationService().DeleteGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, opts...)
}

type MessageCommand struct {
	discord.MessageCommand
	Bot *Bot
}

// Guild returns the Guild the ApplicationCommand is from the Caches or nil if it is a global ApplicationCommand
func (c *MessageCommand) Guild() *Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Bot.Caches.GuildCache().Get(*c.GuildID)
}

// IsGlobal returns true if this is a global ApplicationCommand and false for a guild ApplicationCommand
func (c *MessageCommand) IsGlobal() bool {
	return c.GuildID == nil
}

// ToCreate return the ApplicationCommandCreate for this ApplicationCommand
func (c *MessageCommand) ToCreate() discord.MessageCommandCreate {
	return discord.MessageCommandCreate{
		Name:              c.Name,
		DefaultPermission: c.DefaultPermission,
	}
}

// Update updates the current ApplicationCommand with the given fields
func (c *MessageCommand) Update(commandUpdate discord.MessageCommandUpdate, opts ...rest.RequestOpt) (*MessageCommand, error) {
	var command *discord.ApplicationCommand
	var err error
	if c.GuildID == nil {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGlobalCommand(c.Bot.ApplicationID, c.ID, commandUpdate, opts...)

	} else {
		command, err = c.Bot.RestServices.ApplicationService().UpdateGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, commandUpdate, opts...)
	}
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessageCommand((*command).(discord.MessageCommand)), nil
}

// SetPermissions sets the ApplicationCommandPermissions for a specific Guild. this overrides all existing ApplicationCommandPermission(s). thx discord for that
func (c *MessageCommand) SetPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().SetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// GetPermissions fetched the ApplicationCommandPermissions for a specific Guild from discord
func (c *MessageCommand) GetPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := c.Bot.RestServices.ApplicationService().GetGuildCommandPermissions(c.Bot.ApplicationID, guildID, c.ID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// Delete deletes the ApplicationCommand from discord
func (c *MessageCommand) Delete(opts ...rest.RequestOpt) error {
	if c.GuildID == nil {
		return c.Bot.RestServices.ApplicationService().DeleteGlobalCommand(c.Bot.ApplicationID, c.ID)
	}
	return c.Bot.RestServices.ApplicationService().DeleteGuildCommand(c.Bot.ApplicationID, *c.GuildID, c.ID, opts...)
}
