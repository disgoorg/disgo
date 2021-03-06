package disgo

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events"
)

// Disgo is the main discord interface
type Disgo interface {
	Connect() error
	Close()
	Token() string
	Gateway() Gateway
	RestClient() RestClient
	Intents() models.Intent
	SelfUser() *models.User
	setSelfUser(models.User)
	event(events.GenericEvent)
	AddEventHandlers(...func(event events.GenericEvent))
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token        string
	gateway      Gateway
	restClient   RestClient
	intents      models.Intent
	selfUser     *models.User
	eventHandler EventHandler
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
func (d DisgoImpl) Gateway() Gateway {
	return d.gateway
}

// RestClient returns the HTTP client used by disgo
func (d DisgoImpl) RestClient() RestClient {
	return d.restClient
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

func (d DisgoImpl) setSelfUser(user models.User) {
	d.selfUser = &user
}

func (d DisgoImpl) event(event events.GenericEvent) {
	d.eventHandler.event(event)
}

func (d DisgoImpl) AddEventHandlers(handlers ...func(event events.GenericEvent)) {
	d.eventHandler.AddEventHandlers(handlers...)
}

// New initialises a new disgo client
func New(token string, options Options) Disgo {
	disgo := &DisgoImpl{
		token:        token,
		intents:      options.Intents,
		eventHandler: newEventHandler(),
	}

	disgo.restClient = RestClientImpl{
		disgo:  disgo,
		client: &http.Client{},
	}

	disgo.gateway = &GatewayImpl{
		disgo: disgo,
	}

	return disgo
}
