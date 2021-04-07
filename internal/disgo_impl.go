package internal

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
)

// New creates a new api.Disgo instance
func New(token endpoints.Token, options api.Options) (api.Disgo, error) {
	if options.LargeThreshold < 50 {
		options.LargeThreshold = 50
	} else if options.LargeThreshold > 250 {
		options.LargeThreshold = 250
	}
	disgo := &DisgoImpl{
		BotToken:       token,
		intents:        options.Intents,
		largeThreshold: options.LargeThreshold,
	}

	id, err := IDFromToken(token)
	if err != nil {
		log.Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}

	disgo.selfUserID = *id

	disgo.restClient = newRestClientImpl(disgo)

	disgo.eventManager = newEventManagerImpl(disgo, []api.EventListener{})

	if options.EnableWebhookInteractions {
		disgo.webhookServer = newWebhookServerImpl(disgo, options.ListenURL, options.ListenPort, options.PublicKey)
	}

	disgo.gateway = newGatewayImpl(disgo)

	return disgo, nil
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	// make this public so it does not print in fmt.Sprint("%+v, DisgoImpl{})
	BotToken                 endpoints.Token
	gateway                  api.Gateway
	restClient               api.RestClient
	intents                  api.Intents
	eventManager             api.EventManager
	voiceDispatchInterceptor api.VoiceDispatchInterceptor
	webhookServer            api.WebhookServer
	cache                    api.Cache
	selfUserID               api.Snowflake
	largeThreshold           int
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
func (d *DisgoImpl) Start() {
	d.WebhookServer().Start()
}

// Close will cleanup all disgo internals and close the discord connection safely
func (d *DisgoImpl) Close() {
	if d.RestClient() != nil {
		d.RestClient().Close()
	}
	if d.WebhookServer() != nil {
		d.WebhookServer().Close()
	}
	if d.Gateway() != nil {
		d.Gateway().Close()
	}
	if d.EventManager() != nil {
		d.EventManager().Close()
	}
	if d.Cache() != nil {
		d.Cache().Close()
	}
}

// Token returns the BotToken of the client
func (d *DisgoImpl) Token() endpoints.Token {
	return d.BotToken
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

// SelfUserID returns the current application id
func (d *DisgoImpl) SelfUserID() api.Snowflake {
	return d.selfUserID
}

// SelfUser returns a user object for the client, if available
func (d *DisgoImpl) SelfUser() *api.User {
	return d.cache.User(d.selfUserID)
}

// HeartbeatLatency returns the heartbeat latency
func (d *DisgoImpl) HeartbeatLatency() time.Duration {
	return d.Gateway().Latency()
}

// LargeThreshold returns the large threshold set
func (d *DisgoImpl) LargeThreshold() int {
	return d.largeThreshold
}

// GetCommand fetches a specific guild command
func (d DisgoImpl) GetCommand(commandID api.Snowflake) (*api.Command, error) {
	return d.RestClient().GetGlobalCommand(d.SelfUserID(), commandID)
}

// GetCommands fetches all guild commands
func (d DisgoImpl) GetCommands() ([]*api.Command, error) {
	return d.RestClient().GetGlobalCommands(d.SelfUserID())
}

// CreateCommand creates a new command for this guild
func (d DisgoImpl) CreateCommand(command api.Command) (*api.Command, error) {
	return d.RestClient().CreateGlobalCommand(d.SelfUserID(), command)
}

// EditCommand edits a specific guild command
func (d DisgoImpl) EditCommand(commandID api.Snowflake, command api.Command) (*api.Command, error) {
	return d.RestClient().EditGlobalCommand(d.SelfUserID(), commandID, command)
}

// DeleteCommand creates a new command for this guild
func (d DisgoImpl) DeleteCommand(command api.Command) (*api.Command, error) {
	return d.RestClient().CreateGlobalCommand(d.SelfUserID(), command)
}

// SetCommands overrides all commands for this guild
func (d DisgoImpl) SetCommands(commands ...api.Command) ([]*api.Command, error) {
	return d.RestClient().SetGlobalCommands(d.SelfUserID(), commands...)
}
