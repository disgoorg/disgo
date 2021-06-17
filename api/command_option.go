package api

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
	CommandOptionTypeMentionable
)

// NewCommandOption creates a new CommandOption with the provided params
func NewCommandOption(optionType CommandOptionType, name string, description string, options ...CommandOption) CommandOption {
	return CommandOption{
		Type:        optionType,
		Name:        name,
		Description: description,
		Options:     options,
	}
}

// NewSubCommand creates a new CommandOption with CommandOptionTypeSubCommand
func NewSubCommand(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeSubCommand, name, description, options...)
}

// NewSubCommandGroup creates a new CommandOption with CommandOptionTypeSubCommandGroup
func NewSubCommandGroup(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeSubCommandGroup, name, description, options...)
}

// NewStringOption creates a new CommandOption with CommandOptionTypeSubCommand
func NewStringOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeString, name, description, options...)
}

// NewIntegerOption creates a new CommandOption with CommandOptionTypeSubCommand
func NewIntegerOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeInteger, name, description, options...)
}

// NewBooleanOption creates a new CommandOption with CommandOptionTypeSubCommand
func NewBooleanOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeBoolean, name, description, options...)
}

// NewUserOption creates a new CommandOption with CommandOptionTypeSubCommand
func NewUserOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeUser, name, description, options...)
}

// NewChannelOption creates a new CommandOption with CommandOptionTypeSubCommand
func NewChannelOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeChannel, name, description, options...)
}

// NewRoleOption creates a new CommandOption with CommandOptionTypeRole
func NewRoleOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeRole, name, description, options...)
}

// NewMentionableOption creates a new CommandOption with CommandOptionTypeUser or CommandOptionTypeRole
func NewMentionableOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeMentionable, name, description, options...)
}

// CommandOption are the arguments used in Command.Options
type CommandOption struct {
	Type        CommandOptionType `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required,omitempty"`
	Choices     []OptionChoice    `json:"choices,omitempty"`
	Options     []CommandOption   `json:"options,omitempty"`
}

// AddChoice adds a new choice to the the CommandOption
func (o CommandOption) AddChoice(name string, value interface{}) CommandOption {
	o.Choices = append(o.Choices, OptionChoice{
		Name:  name,
		Value: value,
	})
	return o
}

// AddOptions adds multiple choices to the the CommandOption
func (o CommandOption) AddOptions(options ...CommandOption) CommandOption {
	o.Options = append(o.Options, options...)
	return o
}

// SetRequired sets if the CommandOption is required
func (o CommandOption) SetRequired(required bool) CommandOption {
	o.Required = required
	return o
}

// OptionChoice contains the data for a user using your command
type OptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
