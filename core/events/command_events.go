package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type ApplicationCommandEvent struct {
	*GenericInteractionEvent
	ApplicationCommandInteraction *core.ApplicationCommandInteraction
}

// CommandType returns the type of core.ApplicationCommand which was used
func (e *ApplicationCommandEvent) CommandType() discord.ApplicationCommandType {
	return e.ApplicationCommandInteraction.CommandType()
}

// CommandID returns the ID of the api.Command which got used
func (e *ApplicationCommandEvent) CommandID() discord.Snowflake {
	return e.ApplicationCommandInteraction.CommandID()
}

// CommandName the name of the api.Command which got used
func (e *ApplicationCommandEvent) CommandName() string {
	return e.ApplicationCommandInteraction.CommandName()
}

// Resolved returns the core.Resolved mentions from this core.ApplicationCommand
func (e *ApplicationCommandEvent) Resolved() *core.Resolved {
	return e.ApplicationCommandInteraction.Resolved()
}

// SlashCommandEvent indicates that a slash core.ApplicationCommand was run
type SlashCommandEvent struct {
	*ApplicationCommandEvent
	SlashCommandInteraction *core.SlashCommandInteraction
}

// SubCommandName the subcommand name of the core.ApplicationCommand which got used. May be nil
func (e *SlashCommandEvent) SubCommandName() *string {
	return e.SlashCommandInteraction.SubCommandName()
}

// SubCommandGroupName the subcommand group name of the core.ApplicationCommand which got used. May be nil
func (e *SlashCommandEvent) SubCommandGroupName() *string {
	return e.SlashCommandInteraction.SubCommandGroupName()
}

// CommandPath returns the api.Command path
func (e *SlashCommandEvent) CommandPath() string {
	return e.SlashCommandInteraction.CommandPath()
}

// Options returns the parsed core.ApplicationCommandOption which the core.ApplicationCommand got used with
func (e *SlashCommandEvent) Options() []core.ApplicationCommandOption {
	return e.SlashCommandInteraction.Options()
}

// Option returns an Option by name
func (e *SlashCommandEvent) Option(name string) *core.ApplicationCommandOption {
	return e.SlashCommandInteraction.Option(name)
}

// OptionN returns Option(s) by name
func (e *SlashCommandEvent) OptionN(name string) []core.ApplicationCommandOption {
	return e.SlashCommandInteraction.OptionN(name)
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e *SlashCommandEvent) OptionsT(optionType discord.ApplicationCommandOptionType) []core.ApplicationCommandOption {
	return e.SlashCommandInteraction.OptionsT(optionType)
}

type ContextCommandEvent struct {
	*ApplicationCommandEvent
	ContextCommandInteraction *core.ContextCommandInteraction
}

func (i *ContextCommandEvent) TargetID() discord.Snowflake {
	return i.ContextCommandInteraction.TargetID()
}

type UserCommandEvent struct {
	*ContextCommandEvent
	UserCommandInteraction *core.UserCommandInteraction
}

func (i *UserCommandEvent) TargetUser() *core.User {
	return i.UserCommandInteraction.TargetUser()
}

func (i *UserCommandEvent) TargetMember() *core.Member {
	return i.UserCommandInteraction.TargetMember()
}

type MessageCommandEvent struct {
	*ContextCommandEvent
	MessageCommandInteraction *core.MessageCommandInteraction
}

func (i *MessageCommandEvent) TargetMessage() *core.Message {
	return i.MessageCommandInteraction.TargetMessage()
}
