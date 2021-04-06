package api

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"
)

// Disgo is the main discord interface
type Disgo interface {
	Connect() error
	Start()
	Close()
	Token() string
	Gateway() Gateway
	RestClient() RestClient
	WebhookServer() WebhookServer
	Cache() Cache
	Intents() Intents
	SelfUserID() Snowflake
	SelfUser() *User
	EventManager() EventManager
	VoiceDispatchInterceptor() VoiceDispatchInterceptor
	SetVoiceDispatchInterceptor(voiceInterceptor VoiceDispatchInterceptor)
	HeartbeatLatency() time.Duration
	LargeThreshold() int

	GetCommand(commandID Snowflake) (*SlashCommand, error)
	GetCommands() ([]*SlashCommand, error)
	CreateCommand(command SlashCommand) (*SlashCommand, error)
	EditCommand(commandID Snowflake, command SlashCommand) (*SlashCommand, error)
	DeleteCommand(command SlashCommand) (*SlashCommand, error)
	SetCommands(commands ...SlashCommand) ([]*SlashCommand, error)
}

// EventHandler provides info about the EventHandler
type EventHandler interface {
	Event() GatewayEventName
	New() interface{}
}

// GatewayEventHandler is used to handle raw gateway events
type GatewayEventHandler interface {
	EventHandler
	Handle(Disgo, EventManager, interface{})
}

// WebhookEventHandler is used to handle raw webhook events
type WebhookEventHandler interface {
	EventHandler
	Handle(Disgo, EventManager, chan interface{}, interface{})
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(interface{})
}

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Close()
	AddEventListeners(...EventListener)
	Handle(GatewayEventName, json.RawMessage, chan interface{})
	Dispatch(Event)
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
