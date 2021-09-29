package core

import "github.com/DisgoOrg/disgo/discord"

// NewCommandOption creates a new ApplicationCommandOption with the provided params
func NewCommandOption(optionType discord.ApplicationCommandOptionType, name string, description string, options ...discord.ApplicationCommandOption) *ApplicationCommandOptionBuilder {
	return &ApplicationCommandOptionBuilder{
		ApplicationCommandOption: discord.ApplicationCommandOption{
			Type:        optionType,
			Name:        name,
			Description: description,
			Options:     options,
		},
	}
}

// NewSubCommand creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeSubCommand type
//goland:noinspection GoUnusedExportedFunction
func NewSubCommand(name string, description string, options ...discord.ApplicationCommandOption) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeSubCommand, name, description, options...)
}

// NewSubCommandGroup creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeSubCommandGroup type
//goland:noinspection GoUnusedExportedFunction
func NewSubCommandGroup(name string, description string, options ...discord.ApplicationCommandOption) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeSubCommandGroup, name, description, options...)
}

// NewStringOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeString type
//goland:noinspection GoUnusedExportedFunction
func NewStringOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeString, name, description)
}

// NewIntegerOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeInteger type
//goland:noinspection GoUnusedExportedFunction
func NewIntegerOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeInteger, name, description)
}

// NewBooleanOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeBoolean type
//goland:noinspection GoUnusedExportedFunction
func NewBooleanOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeBoolean, name, description)
}

// NewUserOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeUser type
//goland:noinspection GoUnusedExportedFunction
func NewUserOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeUser, name, description)
}

// NewChannelOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeChannel type
//goland:noinspection GoUnusedExportedFunction
func NewChannelOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeChannel, name, description)
}

// NewRoleOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeRole type
//goland:noinspection GoUnusedExportedFunction
func NewRoleOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeRole, name, description)
}

// NewMentionableOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeMentionable type.
// Such option can accept types discord.ApplicationCommandOptionTypeUser and discord.ApplicationCommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewMentionableOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeMentionable, name, description)
}

// NewNumberOption creates a new ApplicationCommandOption with the discord.ApplicationCommandOptionTypeNumber type
//goland:noinspection GoUnusedExportedFunction
func NewNumberOption(name string, description string) *ApplicationCommandOptionBuilder {
	return NewCommandOption(discord.ApplicationCommandOptionTypeNumber, name, description)
}

type ApplicationCommandOptionBuilder struct {
	discord.ApplicationCommandOption
}

// AddChoice adds a new choice for the ApplicationCommandOption. Value can be a string, int or float
func (b *ApplicationCommandOptionBuilder) AddChoice(name string, value interface{}) *ApplicationCommandOptionBuilder {
	b.Choices = append(b.Choices, discord.ApplicationCommandOptionChoice{
		Name:  name,
		Value: value,
	})
	return b
}

// AddChoices adds multiple choices for the ApplicationCommandOption. Value can be a string, int or float
func (b *ApplicationCommandOptionBuilder) AddChoices(choices map[string]interface{}) *ApplicationCommandOptionBuilder {
	for name, value := range choices {
		b.Choices = append(b.Choices, discord.ApplicationCommandOptionChoice{
			Name:  name,
			Value: value,
		})
	}
	return b
}

// AddOption adds a new ApplicationCommandOption
func (b *ApplicationCommandOptionBuilder) AddOption(optionType discord.ApplicationCommandOptionType, name string, description string) *ApplicationCommandOptionBuilder {
	b.Options = append(b.Options, discord.ApplicationCommandOption{
		Type:        optionType,
		Name:        name,
		Description: description,
	})
	return b
}

// AddOptions adds multiple ApplicationCommandOption
func (b *ApplicationCommandOptionBuilder) AddOptions(options ...discord.ApplicationCommandOption) *ApplicationCommandOptionBuilder {
	b.Options = append(b.Options, options...)
	return b
}

// SetRequired sets if the ApplicationCommandOption is required
func (b *ApplicationCommandOptionBuilder) SetRequired(required bool) *ApplicationCommandOptionBuilder {
	b.Required = required
	return b
}

// Build builds the ApplicationCommandOptionBuilder, returning a discord.ApplicationCommandOption object
func (b *ApplicationCommandOptionBuilder) Build() discord.ApplicationCommandOption {
	return b.ApplicationCommandOption
}
