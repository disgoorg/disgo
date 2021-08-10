package core

import "github.com/DisgoOrg/disgo/discord"

func NewApplicationCommand(commandType discord.ApplicationCommandType, name string, description string, options ...discord.ApplicationCommandOption) *ApplicationCommandBuilder {
	return &ApplicationCommandBuilder{
		ApplicationCommandCreate: discord.ApplicationCommandCreate{
			Type:              commandType,
			Name:              name,
			Description:       description,
			Options:           options,
			DefaultPermission: true,
		},
	}
}

func NewSlashCommand(name string, description string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeSlash, name, description)
}

func NewUserCommand(name string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeUser, name, "")
}

func NewMessageCommand(name string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeMessage, name, "")
}

type ApplicationCommandBuilder struct {
	discord.ApplicationCommandCreate
}

func (b *ApplicationCommandBuilder) SetDefaultPermission(defaultPermission bool) *ApplicationCommandBuilder {
	b.DefaultPermission = defaultPermission
	return b
}

func (b *ApplicationCommandBuilder) Build() discord.ApplicationCommandCreate {
	return b.ApplicationCommandCreate
}
