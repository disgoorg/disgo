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
	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListeners(eventListeners ...EventListener)
	RawEventsEnabled() bool
	VoiceDispatchInterceptor() VoiceDispatchInterceptor
	SetVoiceDispatchInterceptor(voiceInterceptor VoiceDispatchInterceptor)

	HTTPClient() rest.Client
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

	GetCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	GetCommands(opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error)
	CreateCommand(command discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	EditCommand(commandID discord.Snowflake, command discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	DeleteCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error
	SetCommands(commands []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error)

	GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	GetGuildCommands(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error)
	CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	EditGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error)
	DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error
	SetGuildCommands(guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error)

	GetGuildCommandsPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error)
	GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error)
	SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions []discord.GuildCommandPermissionsSet, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error)
	SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions []discord.CommandPermission, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error)

	GetTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, rest.Error)

	GetInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error)
	DeleteInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error)
}
