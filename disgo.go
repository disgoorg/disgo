package disgo

import (
	"net/http"
	"time"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/internal"
)

func New(token string, options api.Options) api.Disgo {
	disgoClient := internal.New(token, options)

	disgoClient.SetRestClient(internal.RestClientImpl{
		DisgoClient:  disgoClient,
		Client: &http.Client{Timeout: time.Duration(options.RestTimeout)},
	})

	eventManager := internal.EventManagerImpl{
		DisgoClient: disgoClient,
		Channel:     make(chan api.GenericEvent),
		Listeners:   &[]*api.EventListener{},
		Handlers:    internal.GetHandlers(),
	}
	disgoClient.SetEventManager(eventManager)

	go eventManager.ListenEvents()

	disgoClient.SetGateway(internal.GatewayImpl{
		DisgoClient: disgoClient,
	})

	return disgoClient
}