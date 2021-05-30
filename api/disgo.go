package api

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/DisgoOrg/log"
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
	SelfUser() *User
	EntityBuilder() EntityBuilder
	EventManager() EventManager
	VoiceDispatchInterceptor() VoiceDispatchInterceptor
	SetVoiceDispatchInterceptor(voiceInterceptor VoiceDispatchInterceptor)
	AudioController() AudioController
	HeartbeatLatency() time.Duration
	LargeThreshold() int
	HasGateway() bool

	GetCommand(commandID Snowflake) (*Command, error)
	GetCommands() ([]*Command, error)
	CreateCommand(command *CommandCreate) (*Command, error)
	EditCommand(commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteCommand(commandID Snowflake) error
	SetCommands(commands ...*CommandCreate) ([]*Command, error)

	GetGuildCommand(guildId Snowflake, commandID Snowflake) (*Command, error)
	GetGuildCommands(guildId Snowflake, ) ([]*Command, error)
	CreateGuildCommand(guildId Snowflake, command *CommandCreate) (*Command, error)
	EditGuildCommand(guildId Snowflake, commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteGuildCommand(guildId Snowflake, commandID Snowflake) error
	SetGuildCommands(guildId Snowflake, commands ...*CommandCreate) ([]*Command, error)

	GetGuildCommandsPermissions(guildId Snowflake) ([]*GuildCommandPermissions, error)
	GetGuildCommandPermissions(guildId Snowflake, commandID Snowflake) (*GuildCommandPermissions, error)
	SetGuildCommandsPermissions(guildId Snowflake, commandPermissions ...*SetGuildCommandPermissions) ([]*GuildCommandPermissions, error)
	SetGuildCommandPermissions(guildId Snowflake, commandID Snowflake, permissions *SetGuildCommandPermissions) (*GuildCommandPermissions, error)
}

// EventHandler provides info about the EventHandler
type EventHandler interface {
	Event() GatewayEventType
	New() interface{}
}

// GatewayEventHandler is used to handle raw gateway events
type GatewayEventHandler interface {
	EventHandler
	HandleGatewayEvent(disgo Disgo, eventManager EventManager, sequenceNumber int, payload interface{})
}

// WebhookEventHandler is used to handle raw webhook events
type WebhookEventHandler interface {
	EventHandler
	HandleWebhookEvent(disgo Disgo, eventManager EventManager, replyChannel chan *InteractionResponse, payload interface{})
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event interface{})
}

// Event the basic interface each event implement
type Event interface {
	Disgo() Disgo
	SequenceNumber() int
}

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Disgo() Disgo
	Close()
	AddEventListeners(eventListeners ...EventListener)
	Handle(eventType GatewayEventType, replyChannel chan *InteractionResponse, sequenceNumber int, payload json.RawMessage)
	Dispatch(event Event)
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
