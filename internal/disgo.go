package internal

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/models"
)

// Options is the configuration used when creating the client
type Options struct {
	Intents     models.Intent
	RestTimeout int
}

func New(token string, options Options) DisgoImpl {
	return DisgoImpl{
		token:        token,
		intents:      options.Intents,
	}
}

func (d DisgoImpl) SetRestClient(restClient RestClientImpl) {
	d.restClient = restClient
}

func (d DisgoImpl) SetEventManager(eventManager EventManagerImpl) {
	d.eventManager = eventManager
}

func (d DisgoImpl) SetGateway(gateway GatewayImpl) {
	d.gateway = gateway
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token        string
	gateway      api.Gateway
	restClient   api.RestClient
	intents      models.Intent
	selfUser     *models.User
	eventManager api.EventManager
	cache        api.Cache
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
	return d.gateway
}

// RestClient returns the HTTP client used by disgo
func (d DisgoImpl) RestClient() api.RestClient {
	return d.restClient
}

// Cache returns the entity cache used by disgo
func (d DisgoImpl) Cache() api.Cache {
	return d.cache
}

// Intents returns the Intents originally specified when creating the client
func (d DisgoImpl) Intents() models.Intent {
	// Todo: Return copy of intents in this method so they can't be modified
	return d.intents
}

// SelfUser returns a user object for the client, if available
func (d DisgoImpl) SelfUser() *models.User {
	return d.selfUser
}

func (d DisgoImpl) SetSelfUser(user models.User) {
	d.selfUser = &user
}

func (d DisgoImpl) EventManager() api.EventManager {
	return d.eventManager
}