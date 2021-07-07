package internal

import (
	"time"

	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/restclient"

	"github.com/DisgoOrg/disgo/api"
)

// New creates a new api.Disgo instance
func New(token string, options api.Options) (api.Disgo, error) {
	if options.LargeThreshold < 50 {
		options.LargeThreshold = 50
	} else if options.LargeThreshold > 250 {
		options.LargeThreshold = 250
	}

	disgo := &DisgoImpl{
		botToken:                token,
		gatewayIntents:          options.GatewayIntents,
		largeThreshold:          options.LargeThreshold,
		logger:                  options.Logger,
		rawGatewayEventsEnabled: options.RawGatewayEventsEnabled,
	}

	id, err := IDFromToken(token)
	if err != nil {
		disgo.Logger().Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}

	disgo.applicationID = *id
	disgo.clientID = *id

	disgo.restClient = newRestClientImpl(disgo, options.HTTPClient)

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
	botToken                 string
	logger                   log.Logger
	gateway                  api.Gateway
	restClient               api.RestClient
	gatewayIntents           api.GatewayIntents
	rawGatewayEventsEnabled  bool
	entityBuilder            api.EntityBuilder
	eventManager             api.EventManager
	voiceDispatchInterceptor api.VoiceDispatchInterceptor
	audioController          api.AudioController
	webhookServer            api.WebhookServer
	cache                    api.Cache
	applicationID            api.Snowflake
	clientID                 api.Snowflake
	selfUser                 *api.SelfUser
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
func (d *DisgoImpl) Token() string {
	return d.botToken
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

// GatewayIntents returns the GatewayIntents originally specified when creating the client
func (d *DisgoImpl) GatewayIntents() api.GatewayIntents {
	// clones the GatewayIntents so they can't be modified
	c := d.gatewayIntents
	return c
}

// RawGatewayEventsEnabled returns if the events.RawGatewayEvent is enabled/disabled
func (d *DisgoImpl) RawGatewayEventsEnabled() bool {
	return d.rawGatewayEventsEnabled
}

// ApplicationID returns the current application id
func (d *DisgoImpl) ApplicationID() api.Snowflake {
	return d.applicationID
}

// SelfUser returns a api.SelfUser for the client, if available
func (d *DisgoImpl) SelfUser() *api.SelfUser {
	return d.selfUser
}

// ClientID returns the current user id
func (d *DisgoImpl) ClientID() api.Snowflake {
	return d.clientID
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
func (d *DisgoImpl) HasGateway() bool {
	return d.gateway != nil
}

// GetCommand fetches a specific global api.Command
func (d *DisgoImpl) GetCommand(commandID api.Snowflake) (*api.Command, restclient.RestError) {
	return d.RestClient().GetGlobalCommand(d.ApplicationID(), commandID)
}

// GetCommands fetches all global api.Command(s)
func (d *DisgoImpl) GetCommands() ([]*api.Command, restclient.RestError) {
	return d.RestClient().GetGlobalCommands(d.ApplicationID())
}

// CreateCommand creates a new global api.Command
func (d *DisgoImpl) CreateCommand(command api.CommandCreate) (*api.Command, restclient.RestError) {
	return d.RestClient().CreateGlobalCommand(d.ApplicationID(), command)
}

// EditCommand edits a specific global api.Command
func (d *DisgoImpl) EditCommand(commandID api.Snowflake, command api.CommandUpdate) (*api.Command, restclient.RestError) {
	return d.RestClient().UpdateGlobalCommand(d.ApplicationID(), commandID, command)
}

// DeleteCommand creates a new global api.Command
func (d *DisgoImpl) DeleteCommand(commandID api.Snowflake) restclient.RestError {
	return d.RestClient().DeleteGlobalCommand(d.ApplicationID(), commandID)
}

// SetCommands overrides all global api.Command(s)
func (d *DisgoImpl) SetCommands(commands ...api.CommandCreate) ([]*api.Command, restclient.RestError) {
	return d.RestClient().SetGlobalCommands(d.ApplicationID(), commands...)
}

// GetGuildCommand fetches a specific api.Guild api.Command
func (d *DisgoImpl) GetGuildCommand(guildID api.Snowflake, commandID api.Snowflake) (*api.Command, restclient.RestError) {
	return d.RestClient().GetGuildCommand(d.ApplicationID(), guildID, commandID)
}

// GetGuildCommands fetches all api.Guild api.Command(s)
func (d *DisgoImpl) GetGuildCommands(guildID api.Snowflake) ([]*api.Command, restclient.RestError) {
	return d.RestClient().GetGuildCommands(d.ApplicationID(), guildID)
}

// CreateGuildCommand creates a new api.Command for this api.Guild
func (d *DisgoImpl) CreateGuildCommand(guildID api.Snowflake, command api.CommandCreate) (*api.Command, restclient.RestError) {
	return d.RestClient().CreateGuildCommand(d.ApplicationID(), guildID, command)
}

// EditGuildCommand edits a specific api.Guild api.Command
func (d *DisgoImpl) EditGuildCommand(guildID api.Snowflake, commandID api.Snowflake, command api.CommandUpdate) (*api.Command, restclient.RestError) {
	return d.RestClient().UpdateGuildCommand(d.ApplicationID(), guildID, commandID, command)
}

// DeleteGuildCommand creates a new api.Command for this api.Guild
func (d *DisgoImpl) DeleteGuildCommand(guildID api.Snowflake, commandID api.Snowflake) restclient.RestError {
	return d.RestClient().DeleteGuildCommand(d.ApplicationID(), guildID, commandID)
}

// SetGuildCommands overrides all api.Command(s) for this api.Guild
func (d *DisgoImpl) SetGuildCommands(guildID api.Snowflake, commands ...api.CommandCreate) ([]*api.Command, restclient.RestError) {
	return d.RestClient().SetGuildCommands(d.ApplicationID(), guildID, commands...)
}

// GetGuildCommandsPermissions returns the api.GuildCommandPermissions for a all api.Command(s) in a api.Guild
func (d *DisgoImpl) GetGuildCommandsPermissions(guildID api.Snowflake) ([]*api.GuildCommandPermissions, restclient.RestError) {
	return d.RestClient().GetGuildCommandsPermissions(d.ApplicationID(), guildID)
}

// GetGuildCommandPermissions returns the api.GuildCommandPermissions for a specific api.Command in a api.Guild
func (d *DisgoImpl) GetGuildCommandPermissions(guildID api.Snowflake, commandID api.Snowflake) (*api.GuildCommandPermissions, restclient.RestError) {
	return d.RestClient().GetGuildCommandPermissions(d.ApplicationID(), guildID, commandID)
}

// SetGuildCommandsPermissions sets the api.GuildCommandPermissions for a all api.Command(s)
func (d *DisgoImpl) SetGuildCommandsPermissions(guildID api.Snowflake, commandPermissions ...api.SetGuildCommandPermissions) ([]*api.GuildCommandPermissions, restclient.RestError) {
	return d.RestClient().SetGuildCommandsPermissions(d.ApplicationID(), guildID, commandPermissions...)
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific api.Command
func (d *DisgoImpl) SetGuildCommandPermissions(guildID api.Snowflake, commandID api.Snowflake, permissions api.SetGuildCommandPermissions) (*api.GuildCommandPermissions, restclient.RestError) {
	return d.RestClient().SetGuildCommandPermissions(d.ApplicationID(), guildID, commandID, permissions)
}

// GetTemplate gets a api.GuildTemplate by it's code
func (d *DisgoImpl) GetTemplate(code string) (*api.GuildTemplate, restclient.RestError) {
	return d.RestClient().GetGuildTemplate(code)
}

// CreateGuildFromTemplate creates an api.Guild using a api.Template code
func (d *DisgoImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate api.CreateGuildFromTemplate) (*api.Guild, restclient.RestError) {
	return d.RestClient().CreateGuildFromTemplate(templateCode, createGuildFromTemplate)
}
