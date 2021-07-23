package events

import "github.com/DisgoOrg/disgo/api"

type GenericCommandEvent struct {
	*GenericInteractionEvent
	GenericCommandInteraction *api.GenericCommandInteraction
}

// CommandID returns the ID of the api.Command which got used
func (e *GenericCommandEvent) CommandID() api.Snowflake {
	return e.GenericCommandInteraction.CommandID()
}

// CommandName the name of the api.Command which got used
func (e *GenericCommandEvent) CommandName() string {
	return e.GenericCommandInteraction.CommandName()
}

// SlashCommandEvent indicates that a slash api.Command was ran
type SlashCommandEvent struct {
	*GenericCommandEvent
	CommandInteraction *api.SlashCommandInteraction
}

// SubCommandName the subcommand name of the api.Command which got used. May be nil
func (e *SlashCommandEvent) SubCommandName() *string {
	return e.CommandInteraction.SubCommandName()
}

// SubCommandGroupName the subcommand group name of the api.Command which got used. May be nil
func (e *SlashCommandEvent) SubCommandGroupName() *string {
	return e.CommandInteraction.SubCommandGroupName()
}

// CommandPath returns the api.Command path
func (e *SlashCommandEvent) CommandPath() string {
	return e.CommandInteraction.CommandPath()
}

// Options returns the parsed api.Option which the api.Command got used with
func (e *SlashCommandEvent) Options() []api.Option {
	return e.CommandInteraction.Options()
}

// Option returns an Option by name
func (e *SlashCommandEvent) Option(name string) *api.Option {
	options := e.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return &options[0]
}

// OptionN returns Option(s) by name
func (e *SlashCommandEvent) OptionN(name string) []api.Option {
	options := make([]api.Option, 0)
	for _, option := range e.Options() {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e *SlashCommandEvent) OptionsT(optionType api.CommandOptionType) []api.Option {
	options := make([]api.Option, 0)
	for _, option := range e.Options() {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

type GenericContextEvent struct {
	*GenericCommandEvent
	GenericContextInteraction *api.GenericContextInteraction
}

func (i *GenericContextEvent) TargetID() api.Snowflake {
	return i.GenericContextInteraction.TargetID()
}

type UserContextEvent struct {
	*GenericContextEvent
	UserContextInteraction *api.UserContextInteraction
}

func (i *UserContextEvent) User() *api.User {
	return i.UserContextInteraction.User()
}

func (i *UserContextEvent) Member() *api.Member {
	return i.UserContextInteraction.Member()
}

type MessageContextEvent struct {
	*GenericContextEvent
	MessageContextInteraction *api.MessageContextInteraction
}

func (i *MessageContextEvent) Message() *api.Message {
	return i.MessageContextInteraction.Message()
}
