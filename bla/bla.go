package bla

import (
	"net/http"
	"time"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/internal"
)

func New(token string, options internal.Options) disgo.Disgo {
	disgoClient := internal.New(token, options)

	disgoClient.SetRestClient(internal.RestClientImpl{
		DisgoClient:  disgoClient,
		Client: &http.Client{Timeout: time.Duration(options.RestTimeout)},
	})

	eventManager := internal.EventManagerImpl{
		DisgoClient: disgoClient,
		Channel:     make(chan disgo.GenericEvent),
		Listeners:   &[]*disgo.EventListener{},
		Handlers:    internal.GetHandlers(),
	}
	go eventManager.ListenEvents()

	disgoClient.SetEventManager(eventManager)

	disgoClient.SetGateway(internal.GatewayImpl{
		DisgoClient: disgoClient,
	})

	return disgoClient
}