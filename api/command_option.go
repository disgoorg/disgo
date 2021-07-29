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
	CommandOptionTypeNumber
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
//goland:noinspection GoUnusedExportedFunction
func NewSubCommand(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeSubCommand, name, description, options...)
}

// NewSubCommandGroup creates a new CommandOption with CommandOptionTypeSubCommandGroup
//goland:noinspection GoUnusedExportedFunction
func NewSubCommandGroup(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeSubCommandGroup, name, description, options...)
}

// NewStringOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewStringOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeString, name, description, options...)
}

// NewIntegerOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewIntegerOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeInteger, name, description, options...)
}

// NewBooleanOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewBooleanOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeBoolean, name, description, options...)
}

// NewUserOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewUserOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeUser, name, description, options...)
}

// NewChannelOption creates a new CommandOption with CommandOptionTypeSubCommand
//goland:noinspection GoUnusedExportedFunction
func NewChannelOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeChannel, name, description, options...)
}

// NewRoleOption creates a new CommandOption with CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewRoleOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeRole, name, description, options...)
}

// NewMentionableOption creates a new CommandOption with CommandOptionTypeUser or CommandOptionTypeRole
//goland:noinspection GoUnusedExportedFunction
func NewMentionableOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeMentionable, name, description, options...)
}

// NewNumberOption creates a new CommandOption with CommandOptionTypeNumber
//goland:noinspection GoUnusedExportedFunction
func NewNumberOption(name string, description string, options ...CommandOption) CommandOption {
	return NewCommandOption(CommandOptionTypeNumber, name, description, options...)
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

// AddChoice adds a new choice to the CommandOption. Value can either be a string, int or float
func (o CommandOption) AddChoice(name string, value interface{}) CommandOption {
	o.Choices = append(o.Choices, OptionChoice{
		Name:  name,
		Value: value,
	})
	return o
}

// AddChoices adds multiple choices to the CommandOption. Value can either be a string, int or float
func (o CommandOption) AddChoices(choices map[string]interface{}) CommandOption {
	for name, value := range choices {
		o.Choices = append(o.Choices, OptionChoice{
			Name:  name,
			Value: value,
		})
	}
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

// OptionChoice contains the data for a user using your command. Value can either be a string, int or float
type OptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
