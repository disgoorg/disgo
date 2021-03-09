package disgo

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/DiscoOrg/disgo/internal"
	"github.com/DiscoOrg/disgo/models"
)


func New(token string, options internal.Options) Disgo {
	disgoClient := internal.New(token, options)

	disgoClient.SetRestClient(internal.RestClientImpl{
		DisgoClient:  disgoClient,
		Client: &http.Client{Timeout: time.Duration(options.RestTimeout)},
	})

	eventManager := internal.EventManagerImpl{
		DisgoClient: disgoClient,
		Channel:     make(chan GenericEvent),
		Listeners:   &[]*EventListener{},
		Handlers:    internal.GetHandlers(),
	}
	go eventManager.ListenEvents()

	disgoClient.SetEventManager(eventManager)

	disgoClient.SetGateway(internal.GatewayImpl{
		DisgoClient: disgoClient,
	})

	return disgoClient
}

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
