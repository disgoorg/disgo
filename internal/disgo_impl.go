package internal

import (
	"time"

	"github.com/DisgoOrg/log"

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
		BotToken:                token,
		intents:                 options.Intents,
		largeThreshold:          options.LargeThreshold,
		logger:                  options.Logger,
		rawGatewayEventsEnabled: options.RawGatewayEventsEnabled,
	}

	id, err := IDFromToken(token)
	if err != nil {
		disgo.Logger().Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}

	disgo.selfUserID = *id

	disgo.restClient = newRestClientImpl(disgo)

	disgo.audioController = newAudioControllerImpl(disgo)

	disgo.entityBuilder = newEntityBuilderImpl(disgo)

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
	logger                   log.Logger
	gateway                  api.Gateway
	restClient               api.RestClient
	intents                  api.Intents
	rawGatewayEventsEnabled  bool
	entityBuilder            api.EntityBuilder
	eventManager             api.EventManager
	voiceDispatchInterceptor api.VoiceDispatchInterceptor
	audioController          api.AudioController
	webhookServer            api.WebhookServer
	cache                    api.Cache
	selfUserID               api.Snowflake
	largeThreshold           int
}

// Logger returns the logger instance disgo uses
func (d *DisgoImpl) Logger() log.Logger {
	return d.logger
}

// Connect opens the gateway connection to discord
func (d *DisgoImpl) Connect() error {
	err := d.Gateway().Open()
	if err != nil {
		d.logger.Errorf("Unable to connect to gateway. error: %s", err)
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

// EntityBuilder returns the api.EntityBuilder
func (d *DisgoImpl) EntityBuilder() api.EntityBuilder {
	return d.entityBuilder
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

// AudioController returns the api.AudioController which can be used to connect/reconnect/disconnect to/fom a api.VoiceChannel
func (d *DisgoImpl) AudioController() api.AudioController {
	return d.audioController
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

// RawGatewayEventsEnabled returns if the events.RawGatewayEvent is enabled/disabled
func (d *DisgoImpl) RawGatewayEventsEnabled() bool {
	return d.rawGatewayEventsEnabled
}

// ApplicationID returns the current application id
func (d *DisgoImpl) ApplicationID() api.Snowflake {
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

// HasGateway returns whether api.Disgo has an active api.Gateway connection or not
func (d DisgoImpl) HasGateway() bool {
	return d.gateway != nil
}

// GetCommand fetches a specific guild command
func (d DisgoImpl) GetCommand(commandID api.Snowflake) (*api.Command, error) {
	return d.RestClient().GetGlobalCommand(d.ApplicationID(), commandID)
}

// GetCommands fetches all guild commands
func (d DisgoImpl) GetCommands() ([]*api.Command, error) {
	return d.RestClient().GetGlobalCommands(d.ApplicationID())
}

// CreateCommand creates a new command for this guild
func (d DisgoImpl) CreateCommand(command api.Command) (*api.Command, error) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// EditCommand edits a specific guild command
func (d DisgoImpl) EditCommand(commandID api.Snowflake, command api.UpdateCommand) (*api.Command, error) {
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
