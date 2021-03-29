package internal

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
)

// New creates a new api.Disgo instance
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

	if options.EnableWebhookInteractions {
		disgo.webhookServer = newWebhookServerImpl(disgo, options.ListenURL, options.ListenPort, options.PublicKey)
	}

	disgo.gateway = newGatewayImpl(disgo)

	return disgo, nil
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	token                    string
	gateway                  api.Gateway
	restClient               api.RestClient
	intents                  api.Intents
	selfUser                 *api.User
	eventManager             api.EventManager
	voiceDispatchInterceptor api.VoiceDispatchInterceptor
	webhookServer            api.WebhookServer
	cache                    api.Cache
	applicationID            api.Snowflake
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

// Start starts the interaction webhook server
func (d *DisgoImpl) Start() error {
	err := d.WebhookServer().Start()
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

// EventManager returns the api.EventManager
func (d *DisgoImpl) EventManager() api.EventManager {
	return d.eventManager
}

// VoiceDispatchInterceptor returns the api.VoiceDispatchInterceptor
func (d *DisgoImpl) VoiceDispatchInterceptor() api.VoiceDispatchInterceptor {
	return d.voiceDispatchInterceptor
}

// SetVoiceDispatchInterceptor sets the api.VoiceDispatchInterceptor
func (d *DisgoImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor api.VoiceDispatchInterceptor) {
	d.voiceDispatchInterceptor = voiceDispatchInterceptor
}

// WebhookServer returns the api.EventManager
func (d *DisgoImpl) WebhookServer() api.WebhookServer {
	return d.webhookServer
}

// Cache returns the entity api.Cache used by disgo
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

// SetSelfUser sets the self user
func (d *DisgoImpl) SetSelfUser(user *api.User) {
	d.selfUser = user
}

// HeartbeatLatency returns the heartbeat latency
func (d *DisgoImpl) HeartbeatLatency() time.Duration {
	return d.Gateway().Latency()
}

// GetCommand fetches a specific guild command
func (d DisgoImpl) GetCommand(commandID api.Snowflake) (*api.SlashCommand, error) {
	return d.RestClient().GetGlobalCommand(d.ApplicationID(), commandID)
}

// GetCommands fetches all guild commands
func (d DisgoImpl) GetCommands() ([]*api.SlashCommand, error) {
	return d.RestClient().GetGlobalCommands(d.ApplicationID())
}

// CreateCommand creates a new command for this guild
func (d DisgoImpl) CreateCommand(command api.SlashCommand) (*api.SlashCommand, error) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// EditCommand edits a specific guild command
func (d DisgoImpl) EditCommand(commandID api.Snowflake, command api.SlashCommand) (*api.SlashCommand, error) {
	return d.RestClient().EditGlobalCommand(d.ApplicationID(), commandID, command)
}

// DeleteCommand creates a new command for this guild
func (d DisgoImpl) DeleteCommand(command api.SlashCommand) (*api.SlashCommand, error) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// SetCommands overrides all commands for this guild
func (d DisgoImpl) SetCommands(commands ...api.SlashCommand) ([]*api.SlashCommand, error) {
	return d.RestClient().SetGlobalCommands(d.ApplicationID(), commands...)
}
