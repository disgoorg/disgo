package api

import (
	"encoding/json"
	"runtime"
	"strings"
)

// Disgo is the main discord interface
type Disgo interface {
	Connect() error
	Close()
	Token() string
	Gateway() Gateway
	RestClient() RestClient
	Cache() Cache
	Intents() Intent
	SelfUser() *User
	SetSelfUser(User)
	EventManager() EventManager
}

type GatewayEventProvider interface {
	New() interface{}
	Handle(Disgo, EventManager, interface{})
}

type EventListener interface {
	OnEvent(interface{})
}

type EventManager interface {
	AddEventListeners(...EventListener)
	Handle(string, json.RawMessage)
	Dispatch(GenericEvent)
}

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
