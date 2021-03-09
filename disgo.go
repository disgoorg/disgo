package disgo

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/models"
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
	setSelfUser(models.User)
	EventManager() EventManager
}

// disgoImpl is the main discord client
type disgoImpl struct {
	token        string
	gateway      Gateway
	restClient   RestClient
	intents      models.Intent
	selfUser     *models.User
	eventManager EventManager
	cache        Cache
}

// Connect opens the gateway connection to discord
func (d disgoImpl) Connect() error {
	err := d.Gateway().Open()
	if err != nil {
		log.Errorf("Unable to connect to gateway. error: %s", err)
		return err
	}
	return nil
}

// Close will cleanup all disgo internals and close the discord connection safely
func (d disgoImpl) Close() {
	d.RestClient().Close()
	d.Gateway().Close()
}

// Token returns the token of the client
func (d disgoImpl) Token() string {
	return d.token
}

// Gateway returns the websocket information
func (d disgoImpl) Gateway() Gateway {
	return d.gateway
}

// RestClient returns the HTTP client used by disgo
func (d disgoImpl) RestClient() RestClient {
	return d.restClient
}

// Cache returns the entity cache used by disgo
func (d disgoImpl) Cache() Cache {
	return d.cache
}

// Intents returns the Intents originally specified when creating the client
func (d disgoImpl) Intents() models.Intent {
	// Todo: Return copy of intents in this method so they can't be modified
	return d.intents
}

// SelfUser returns a user object for the client, if available
func (d disgoImpl) SelfUser() *models.User {
	return d.selfUser
}

func (d disgoImpl) setSelfUser(user models.User) {
	d.selfUser = &user
}

func (d disgoImpl) EventManager() EventManager {
	return d.eventManager
}

// New initialises a new disgo client
func New(token string, options Options) Disgo {
	disgo := &disgoImpl{
		token:        token,
		intents:      options.Intents,
	}

	disgo.restClient = RestClientImpl{
		disgo:  disgo,
		client: &http.Client{},
	}

	disgo.eventManager = newEventManager(disgo)

	disgo.gateway = &GatewayImpl{
		disgo: disgo,
	}

	return disgo
}
