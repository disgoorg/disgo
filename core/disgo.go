package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// Disgo is the main discord interface
type Disgo interface {
	Logger() log.Logger
	Close()

	Token() string
	ApplicationID() discord.Snowflake
	ClientID() discord.Snowflake
	SelfUser() *SelfUser
	SetSelfUser(selfUser *SelfUser)
	SelfMember(guildID discord.Snowflake) *Member

	EventManager() EventManager
	RawEventsEnabled() bool
	VoiceDispatchInterceptor() VoiceDispatchInterceptor
	SetVoiceDispatchInterceptor(voiceInterceptor VoiceDispatchInterceptor)

	HTTPClient() rest.HTTPClient
	RestServices() rest.Services

	Gateway() gateway.Gateway
	Connect() error
	HasGateway() bool

	HTTPServer() httpserver.Server
	Start()
	HasHTTPServer() bool

	Cache() Cache

	EntityBuilder() EntityBuilder
	AudioController() AudioController

	GetCommand(commandID discord.Snowflake) (*Command, rest.Error)
	GetCommands() ([]*Command, rest.Error)
	CreateCommand(command discord.CommandCreate) (*Command, rest.Error)
	EditCommand(commandID discord.Snowflake, command discord.CommandUpdate) (*Command, rest.Error)
	DeleteCommand(commandID discord.Snowflake) rest.Error
	SetCommands(commands ...discord.CommandCreate) ([]*Command, rest.Error)

	GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake) (*Command, rest.Error)
	GetGuildCommands(guildID discord.Snowflake) ([]*Command, rest.Error)
	CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.CommandCreate) (*Command, rest.Error)
	EditGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.CommandUpdate) (*Command, rest.Error)
	DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake) rest.Error
	SetGuildCommands(guildID discord.Snowflake, commandCreates ...discord.CommandCreate) ([]*Command, rest.Error)

	GetGuildCommandsPermissions(guildID discord.Snowflake) ([]*GuildCommandPermissions, rest.Error)
	GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake) (*GuildCommandPermissions, rest.Error)
	SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) ([]*GuildCommandPermissions, rest.Error)
	SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions ...discord.CommandPermission) (*GuildCommandPermissions, rest.Error)

	GetTemplate(templateCode string) (*GuildTemplate, rest.Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*Guild, rest.Error)
}
