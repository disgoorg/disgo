package core

import "github.com/DisgoOrg/disgo/discord"

// NewCommandOption creates a new CommandOption with the provided params
func NewCommandOption(optionType discord.CommandOptionType, name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return discord.CommandOption{
		Type:        optionType,
		Name:        name,
		Description: description,
		Options:     options,
	}
}

// NewSubCommand creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewSubCommand(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeSubCommand, name, description, options...)
}

// NewSubCommandGroup creates a new CommandOption with CommandOptionTypeSubCommandGroup
//goland:noinspection GoUnusedExportedFunction
func NewSubCommandGroup(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeSubCommandGroup, name, description, options...)
}

// NewStringOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewStringOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeString, name, description, options...)
}

// NewIntegerOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewIntegerOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeInteger, name, description, options...)
}

// NewBooleanOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewBooleanOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeBoolean, name, description, options...)
}

// NewUserOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewUserOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeUser, name, description, options...)
}

// NewChannelOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewChannelOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeChannel, name, description, options...)
}

// NewRoleOption creates a new CommandOption with CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewRoleOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeRole, name, description, options...)
}

// NewMentionableOption creates a new CommandOption with CommandOptionTypeUser or CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewMentionableOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeMentionable, name, description, options...)
}

// NewNumberOption creates a new CommandOption with CommandOptionTypeNumber
//goland:noinspection GoUnusedExportedFunction
func NewNumberOption(name string, description string, options ...discord.CommandOption) discord.CommandOption {
	return NewCommandOption(discord.CommandOptionTypeNumber, name, description, options...)
}
