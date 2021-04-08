package api

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/DisgoOrg/disgo/api/endpoints"
)

// Disgo is the main discord interface
type Disgo interface {
	Connect() error
	Start()
	Close()
	Token() endpoints.Token
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

	GetCommand(commandID Snowflake) (*Command, error)
	GetCommands() ([]*Command, error)
	CreateCommand(command Command) (*Command, error)
	EditCommand(commandID Snowflake, command UpdateCommand) (*Command, error)
	DeleteCommand(command Command) (*Command, error)
	SetCommands(commands ...Command) ([]*Command, error)
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

// Event the basic interface each event implement
type Event interface {
	Disgo() Disgo
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
