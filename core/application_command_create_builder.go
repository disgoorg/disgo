package core

import "github.com/DisgoOrg/disgo/discord"

func NewApplicationCommand(commandType discord.ApplicationCommandType, name string, description string, options ...discord.SlashCommandOption) *ApplicationCommandBuilder {
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

//goland:noinspection GoUnusedExportedFunction
func NewSlashCommand(name string, description string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeSlash, name, description)
}

//goland:noinspection GoUnusedExportedFunction
func NewUserCommand(name string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeUser, name, "")
}

//goland:noinspection GoUnusedExportedFunction
func NewMessageCommand(name string) *ApplicationCommandBuilder {
	return NewApplicationCommand(discord.ApplicationCommandTypeMessage, name, "")
}

// TODO: complete this?

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
