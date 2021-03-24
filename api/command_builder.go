package api

// NewGuildCommandBuilder creates a new GuildCommandBuilder for creating slash commands
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

// NewGlobalCommandBuilder creates a new GlobalCommandBuilder for creating slash commands
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

// CommandBuilder is a generic builder for creating commands
type CommandBuilder struct {
	disgo   Disgo
	command ApplicationCommand
}

// Build returns the finished ApplicationCommand
func (b CommandBuilder) Build() ApplicationCommand {
	return b.command
}

// GuildCommandBuilder extends CommandBuilder to create guild-specific commands
type GuildCommandBuilder struct {
	CommandBuilder
	guildID Snowflake
}

// Create POSTs your command to discord
func (b GuildCommandBuilder) Create() (*ApplicationCommand, error) {
	return b.disgo.RestClient().CreateGuildApplicationGuildCommand(b.guildID, b.disgo.ApplicationID(), b.command)
}

// GlobalCommandBuilder extends CommandBuilder to create global/DM commands
type GlobalCommandBuilder struct {
	CommandBuilder
}

// Create POSTs your command to discord
func (b GlobalCommandBuilder) Create() (*ApplicationCommand, error) {
	return b.disgo.RestClient().CreateGlobalApplicationGlobalCommand(b.disgo.ApplicationID(), b.command)
}