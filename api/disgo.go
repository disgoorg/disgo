package api

import (
	"runtime"
	"strings"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/restclient"
)

// Disgo is the main discord interface
type Disgo interface {
	Logger() log.Logger
	Connect() error
	Start()
	Close()
	Token() string
	Gateway() Gateway
	RestClient() RestClient
	WebhookServer() WebhookServer
	Cache() Cache
	GatewayIntents() GatewayIntents
	RawGatewayEventsEnabled() bool
	ApplicationID() Snowflake
	ClientID() Snowflake
	SelfUser() *SelfUser
	EntityBuilder() EntityBuilder
	EventManager() EventManager
	VoiceDispatchInterceptor() VoiceDispatchInterceptor
	SetVoiceDispatchInterceptor(voiceInterceptor VoiceDispatchInterceptor)
	AudioController() AudioController
	HeartbeatLatency() time.Duration
	LargeThreshold() int
	HasGateway() bool

	GetCommand(commandID Snowflake) (*Command, restclient.RestError)
	GetCommands() ([]*Command, restclient.RestError)
	CreateCommand(command CommandCreate) (*Command, restclient.RestError)
	EditCommand(commandID Snowflake, command CommandUpdate) (*Command, restclient.RestError)
	DeleteCommand(commandID Snowflake) restclient.RestError
	SetCommands(commands ...CommandCreate) ([]*Command, restclient.RestError)

	GetGuildCommand(guildID Snowflake, commandID Snowflake) (*Command, restclient.RestError)
	GetGuildCommands(guildID Snowflake) ([]*Command, restclient.RestError)
	CreateGuildCommand(guildID Snowflake, commandCreate CommandCreate) (*Command, restclient.RestError)
	EditGuildCommand(guildID Snowflake, commandID Snowflake, commandUpdate CommandUpdate) (*Command, restclient.RestError)
	DeleteGuildCommand(guildID Snowflake, commandID Snowflake) restclient.RestError
	SetGuildCommands(guildID Snowflake, commandCreates ...CommandCreate) ([]*Command, restclient.RestError)

	GetGuildCommandsPermissions(guildID Snowflake) ([]*GuildCommandPermissions, restclient.RestError)
	GetGuildCommandPermissions(guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, restclient.RestError)
	SetGuildCommandsPermissions(guildID Snowflake, commandPermissions ...SetGuildCommandPermissions) ([]*GuildCommandPermissions, restclient.RestError)
	SetGuildCommandPermissions(guildID Snowflake, commandID Snowflake, permissions SetGuildCommandPermissions) (*GuildCommandPermissions, restclient.RestError)

	GetTemplate(code string) (*GuildTemplate, restclient.RestError)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate CreateGuildFromTemplate) (*Guild, restclient.RestError)
}

// GetOS returns the simplified version of the operating system for sending to Discord in the IdentifyCommandDataProperties.OS payload
func GetOS() string {
	OS := runtime.GOOS
	if strings.HasPrefix(OS, "windows") {
		return "windows"
	}
	if strings.HasPrefix(OS, "darwin") {
		return "darwin"
	}
	return "linux"
}
