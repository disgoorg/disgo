package api

import "github.com/chebyrash/promise"

func NewGuildCommandBuilder(disgo Disgo, guildID Snowflake, name string, description string) GuildCommandBuilder {
	return GuildCommandBuilder{
		CommandBuilder: CommandBuilder{
			disgo: disgo,
			command: ApplicationCommand{
				Name: name,
				Description: description,
			},
		},
		guildID: guildID,
	}
}

func NewGlobalCommandBuilder(disgo Disgo, name string, description string) GlobalCommandBuilder {
	return GlobalCommandBuilder{
		CommandBuilder: CommandBuilder{
			disgo: disgo,
			command: ApplicationCommand{
				Name: name,
				Description: description,
			},
		},
	}
}

type CommandBuilder struct {
	disgo   Disgo
	command ApplicationCommand
}

func (b CommandBuilder) Build() ApplicationCommand {
	return b.command
}

type GuildCommandBuilder struct {
	CommandBuilder
	guildID Snowflake
}

func (b GuildCommandBuilder) Create() *promise.Promise {
	return b.disgo.RestClient().CreateGuildApplicationGuildCommand(b.guildID, b.disgo.ApplicationID(), b.command)
}

type GlobalCommandBuilder struct {
	CommandBuilder
}

func (b GlobalCommandBuilder) Create() *promise.Promise {
	return b.disgo.RestClient().CreateGlobalApplicationGlobalCommand(b.disgo.ApplicationID(), b.command)
}