package core

import "github.com/DisgoOrg/disgo/discord"

// NewCommandOption creates a new SlashCommandOption with the provided params
func NewCommandOption(optionType discord.SlashCommandOptionType, name string, description string, options ...discord.SlashCommandOption) *ApplicationCommandOptionBuilder {
	return &ApplicationCommandOptionBuilder{
		SlashCommandOption: discord.SlashCommandOption{
			Type:        optionType,
			Name:        name,
			Description: description,
			Options:     options,
		},
	}
}

// NewSubCommand creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewSubCommand(name string, description string, options ...discord.SlashCommandOption) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeSubCommand, name, description, options...)
}

// NewSubCommandGroup creates a new SlashCommandOption with CommandOptionTypeSubCommandGroup
//goland:noinspection GoUnusedExportedFunction
func NewSubCommandGroup(name string, description string, options ...discord.SlashCommandOption) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeSubCommandGroup, name, description, options...)
}

// NewStringOption creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewStringOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeString, name, description)
}

// NewIntegerOption creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewIntegerOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeInteger, name, description)
}

// NewBooleanOption creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewBooleanOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeBoolean, name, description)
}

// NewUserOption creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewUserOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeUser, name, description)
}

// NewChannelOption creates a new SlashCommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewChannelOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeChannel, name, description)
}

// NewRoleOption creates a new SlashCommandOption with CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewRoleOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeRole, name, description)
}

// NewMentionableOption creates a new SlashCommandOption with CommandOptionTypeUser or CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewMentionableOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeMentionable, name, description)
}

// NewNumberOption creates a new SlashCommandOption with CommandOptionTypeNumber
//goland:noinspection GoUnusedExportedFunction
func NewNumberOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.CommandOptionTypeNumber, name, description)
}

type ApplicationCommandOptionBuilder struct {
	discord.SlashCommandOption
}

// AddChoice adds a new choice to the SlashCommandOption. Value can either be a string, int or float
func (b *ApplicationCommandOptionBuilder) AddChoice(name string, value interface{}) *ApplicationCommandOptionBuilder {
	b.Choices = append(b.Choices, discord.SlashCommandOptionChoice{
		Name:  name,
		Value: value,
	})
	return b
}

// AddChoices adds multiple choices to the SlashCommandOption. Value can either be a string, int or float
func (b *ApplicationCommandOptionBuilder) AddChoices(choices map[string]interface{}) *ApplicationCommandOptionBuilder {
	for name, value := range choices {
		b.Choices = append(b.Choices, discord.SlashCommandOptionChoice{
			Name:  name,
			Value: value,
		})
	}
	return b
}

// AddOption adds a new discord.SlashCommandOption
func (b *ApplicationCommandOptionBuilder) AddOption(optionType discord.SlashCommandOptionType, name string, description string) *ApplicationCommandOptionBuilder {
	b.Options = append(b.Options, discord.SlashCommandOption{
		Type:        optionType,
		Name:        name,
		Description: description,
	})
	return b
}

// AddOptions adds multiple choices to the SlashCommandOption
func (b *ApplicationCommandOptionBuilder) AddOptions(options ...discord.SlashCommandOption) *ApplicationCommandOptionBuilder {
	b.Options = append(b.Options, options...)
	return b
}

// SetRequired sets if the SlashCommandOption is required
func (b *ApplicationCommandOptionBuilder) SetRequired(required bool) *ApplicationCommandOptionBuilder {
	b.Required = required
	return b
}

// Build builds the ApplicationCommandOptionBuilder to discord.SlashCommandOption
func (b *ApplicationCommandOptionBuilder) Build() discord.SlashCommandOption {
	return b.SlashCommandOption
}
