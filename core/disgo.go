package core

import (
	"context"

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

	GetCommand(ctx context.Context, commandID discord.Snowflake) (*ApplicationCommand, rest.Error)
	GetCommands(ctx context.Context) ([]*ApplicationCommand, rest.Error)
	CreateCommand(ctx context.Context, command discord.ApplicationCommandCreate) (*ApplicationCommand, rest.Error)
	EditCommand(ctx context.Context, commandID discord.Snowflake, command discord.ApplicationCommandUpdate) (*ApplicationCommand, rest.Error)
	DeleteCommand(ctx context.Context, commandID discord.Snowflake) rest.Error
	SetCommands(ctx context.Context, commands ...discord.ApplicationCommandCreate) ([]*ApplicationCommand, rest.Error)

	GetGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) (*ApplicationCommand, rest.Error)
	GetGuildCommands(ctx context.Context, guildID discord.Snowflake) ([]*ApplicationCommand, rest.Error)
	CreateGuildCommand(ctx context.Context, guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate) (*ApplicationCommand, rest.Error)
	EditGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (*ApplicationCommand, rest.Error)
	DeleteGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) rest.Error
	SetGuildCommands(ctx context.Context, guildID discord.Snowflake, commandCreates ...discord.ApplicationCommandCreate) ([]*ApplicationCommand, rest.Error)

	GetGuildCommandsPermissions(ctx context.Context, guildID discord.Snowflake) ([]*GuildCommandPermissions, rest.Error)
	GetGuildCommandPermissions(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) (*GuildCommandPermissions, rest.Error)
	SetGuildCommandsPermissions(ctx context.Context, guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) ([]*GuildCommandPermissions, rest.Error)
	SetGuildCommandPermissions(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake, permissions ...discord.CommandPermission) (*GuildCommandPermissions, rest.Error)

	GetTemplate(ctx context.Context, templateCode string) (*GuildTemplate, rest.Error)
	CreateGuildFromTemplate(ctx context.Context, templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*Guild, rest.Error)

	GetInvite(ctx context.Context, inviteCode string) (*Invite, rest.Error)
	DeleteInvite(ctx context.Context, inviteCode string) (*Invite, rest.Error)
}
