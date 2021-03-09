package api

import (
	"encoding/json"
	"runtime"
	"strings"

	"github.com/DiscoOrg/disgo/api/models"
)

// Disgo is the main discord interface
type Disgo interface {
	Connect() error
	Close()
	Token() string
	Gateway() Gateway
	RestClient() RestClient
	Cache() Cache
	Intents() models.Intent
	SelfUser() *models.User
	SetSelfUser(models.User)
	EventManager() EventManager
}

type GatewayEventProvider interface {
	New() interface{}
	Handle(EventManager, interface{})
}

type EventListener interface {
	OnEvent(interface{})
}

type EventManager interface {
	AddEventListeners(...EventListener)
	Handle(string, json.RawMessage)
	Dispatch(GenericEvent)
	Disgo() Disgo
}


// Gateway is what is used to connect to discord
type Gateway interface {
	Disgo() Disgo
	Open() error
	Close()
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
