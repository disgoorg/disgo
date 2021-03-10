package internal

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

func New(token string, options api.Options) DisgoImpl {
	return DisgoImpl{
		token:        token,
		intents:      options.Intents,
	}
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token        string
	gateway      *api.Gateway
	restClient   *api.RestClient
	intents      api.Intent
	selfUser     *api.User
	eventManager *api.EventManager
	cache        api.Cache
}

func (d DisgoImpl) SetRestClient(restClient api.RestClient) {
	d.restClient = &restClient
}

func (d DisgoImpl) SetEventManager(eventManager api.EventManager) {
	d.eventManager = &eventManager
}

func (d DisgoImpl) SetGateway(gateway api.Gateway) {
	d.gateway = &gateway
}

// Connect opens the gateway connection to discord
func (d DisgoImpl) Connect() error {
	err := d.Gateway().Open()
	if err != nil {
		log.Errorf("Unable to connect to gateway. error: %s", err)
		return err
	}
	return nil
}

// Close will cleanup all disgo internals and close the discord connection safely
func (d DisgoImpl) Close() {
	d.RestClient().Close()
	d.Gateway().Close()
}

// Token returns the token of the client
func (d DisgoImpl) Token() string {
	return d.token
}

// Gateway returns the websocket information
func (d DisgoImpl) Gateway() api.Gateway {
	return *d.gateway
}

// RestClient returns the HTTP client used by disgo
func (d DisgoImpl) RestClient() api.RestClient {
	return *d.restClient
}

func (d DisgoImpl) EventManager() api.EventManager {
	return *d.eventManager
}

// Cache returns the entity cache used by disgo
func (d DisgoImpl) Cache() api.Cache {
	return d.cache
}

// Intents returns the Intents originally specified when creating the client
func (d DisgoImpl) Intents() api.Intent {
	// Todo: Return copy of intents in this method so they can't be modified
	return d.intents
}

// SelfUser returns a user object for the client, if available
func (d DisgoImpl) SelfUser() *api.User {
	return d.selfUser
}

func (d DisgoImpl) SetSelfUser(user api.User) {
	d.selfUser = &user
}
