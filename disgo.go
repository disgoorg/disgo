package disgo

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/disgo/models"
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
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token      string
	gateway    Gateway
	restClient RestClient
	intents    models.Intent
	selfUser   *models.User
}

func (d DisgoImpl) Connect() error {
	err := d.Gateway().Open()
	if err != nil {
		log.Errorf("Unable to connect to gateway. error: %s", err)
		return err
	}
	return nil
}

func (d DisgoImpl) Close() {
	d.RestClient().Close()
	d.Gateway().Close()
}

func (d DisgoImpl) Token() string {
	return d.token
}

func (d DisgoImpl) Gateway() Gateway {
	return d.gateway
}

func (d DisgoImpl) RestClient() RestClient {
	return d.restClient
}

func (d DisgoImpl) Intents() models.Intent {
	return d.intents
}

func (d DisgoImpl) SelfUser() *models.User {
	return d.selfUser
}

func (d DisgoImpl) setSelfUser(user models.User) {
	d.selfUser = &user
}

func New(token string, options Options) Disgo {
	disgo := &DisgoImpl{
		token:   token,
		intents: options.Intents,
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
