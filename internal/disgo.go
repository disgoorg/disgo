package internal

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

func New(token string, options api.Options) (api.Disgo, error) {
	disgo := &DisgoImpl{
		token:   token,
		intents: options.Intents,
	}

	id, err := IDFromToken(token)
	if err != nil {
		log.Errorf("error while getting application id from token: %s", err)
		return nil, err
	}

	disgo.applicationID = *id

	disgo.restClient = newRestClientImpl(disgo, token)

	disgo.eventManager = newEventManagerImpl(disgo, make([]api.EventListener, 0))

	disgo.gateway = newGatewayImpl(disgo)

	return disgo, nil
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token         string
	gateway       api.Gateway
	restClient    api.RestClient
	intents       api.Intents
	selfUser      *api.User
	eventManager  api.EventManager
	cache         api.Cache
	applicationID api.Snowflake
}

// Connect opens the gateway connection to discord
func (d *DisgoImpl) Connect() error {
	err := d.Gateway().Open()
	if err != nil {
		log.Errorf("Unable to connect to gateway. error: %s", err)
		return err
	}
	return nil
}

// Close will cleanup all disgo internals and close the discord connection safely
func (d *DisgoImpl) Close() {
	d.RestClient().Close()
	d.Gateway().Close()
}

// Token returns the token of the client
func (d *DisgoImpl) Token() string {
	return d.token
}

// Gateway returns the websocket information
func (d *DisgoImpl) Gateway() api.Gateway {
	return d.gateway
}

// RestClient returns the HTTP client used by disgo
func (d *DisgoImpl) RestClient() api.RestClient {
	return d.restClient
}

func (d *DisgoImpl) EventManager() api.EventManager {
	return d.eventManager
}

// Cache returns the entity cache used by disgo
func (d *DisgoImpl) Cache() api.Cache {
	return d.cache
}

// Intents returns the Intents originally specified when creating the client
func (d *DisgoImpl) Intents() api.Intents {
	// clones the intents so they can't be modified
	c := d.intents
	return c
}

// ApplicationID returns the current application id
func (d *DisgoImpl) ApplicationID() api.Snowflake {
	return d.applicationID
}

// SelfUser returns a user object for the client, if available
func (d *DisgoImpl) SelfUser() *api.User {
	return d.selfUser
}

func (d *DisgoImpl) SetSelfUser(user *api.User) {
	d.selfUser = user
}

func (d *DisgoImpl) HeartbeatLatency() time.Duration {
	return d.Gateway().Latency()
}

// GetCommand fetches a specific guild command
func (d DisgoImpl) GetCommand(commandID api.Snowflake) (*api.Command, error) {
	return d.RestClient().GetGlobalCommand(d.ApplicationID(), commandID)
}

// GetCommand fetches all guild commands
func (d DisgoImpl) GetCommands() ([]*api.Command, error) {
	return d.RestClient().GetGlobalCommands(d.ApplicationID())
}

// CreateCommand creates a new command for this guild
func (d DisgoImpl) CreateCommand(command api.Command) (*api.Command, error) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// EditCommand edits a specific guild command
func (d DisgoImpl) EditCommand(commandID api.Snowflake, command api.Command) (*api.Command, error) {
	return d.RestClient().EditGlobalCommand(d.ApplicationID(), commandID, command)
}

// DeleteCommand creates a new command for this guild
func (d DisgoImpl) DeleteCommand(command api.Command) (*api.Command, error) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// SetCommands overrides all commands for this guild
func (d DisgoImpl) SetCommands(commands ...api.Command) ([]*api.Command, error) {
	return d.RestClient().SetGlobalCommands(d.ApplicationID(), commands...)
}
