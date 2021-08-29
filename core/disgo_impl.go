package core

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest/rate"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

var _ Disgo = (*DisgoImpl)(nil)

func NewDisgo(token string, opts ...DisgoOpt) (Disgo, error) {
	config := &DisgoConfig{Token: token}
	config.Apply(opts)

	return buildDisgoImpl(*config)
}

// DisgoImpl is the main discord client
type DisgoImpl struct {
	logger log.Logger

	token         string
	applicationID discord.Snowflake
	clientID      discord.Snowflake
	selfUser      *SelfUser

	httpClient   rest.Client
	restServices rest.Services

	eventManager             EventManager
	rawEventsEnabled         bool
	voiceDispatchInterceptor VoiceDispatchInterceptor

	gateway gateway.Gateway

	httpServer httpserver.Server

	cache Cache

	entityBuilder   EntityBuilder
	audioController AudioController
}

// Logger returns the logger instance disgo uses
func (d *DisgoImpl) Logger() log.Logger {
	return d.logger
}

// Close will clean up all disgo internals and close the discord connection safely
func (d *DisgoImpl) Close() {
	if d.HTTPClient() != nil {
		d.HTTPClient().Close()
	}
	if d.RestServices() != nil {
		d.RestServices().Close()
	}
	if d.HTTPServer() != nil {
		d.HTTPServer().Close()
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
	return d.token
}

// ApplicationID returns the current application id
func (d *DisgoImpl) ApplicationID() discord.Snowflake {
	return d.applicationID
}

// ClientID returns the current user id
func (d *DisgoImpl) ClientID() discord.Snowflake {
	return d.clientID
}

// SelfUser returns an SelfUser for the client, if available
func (d *DisgoImpl) SelfUser() *SelfUser {
	return d.selfUser
}

// SetSelfUser sets the current SelfUser of Disgo
func (d *DisgoImpl) SetSelfUser(selfUser *SelfUser) {
	d.selfUser = selfUser
}

// SelfMember returns an api.SelfUser for the client, if available
func (d *DisgoImpl) SelfMember(guildID discord.Snowflake) *Member {
	return d.Cache().MemberCache().Get(guildID, d.clientID)
}

// EventManager returns the api.EventManager
func (d *DisgoImpl) EventManager() EventManager {
	return d.eventManager
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (d *DisgoImpl) AddEventListeners(listeners ...EventListener) {
	d.EventManager().AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (d *DisgoImpl) RemoveEventListeners(listeners ...EventListener) {
	d.EventManager().RemoveEventListeners(listeners...)
}

// RawEventsEnabled returns if the events.RawGatewayEvent is enabled/disabled
func (d *DisgoImpl) RawEventsEnabled() bool {
	return d.rawEventsEnabled
}

// VoiceDispatchInterceptor returns the VoiceDispatchInterceptor
func (d *DisgoImpl) VoiceDispatchInterceptor() VoiceDispatchInterceptor {
	return d.voiceDispatchInterceptor
}

// SetVoiceDispatchInterceptor sets the VoiceDispatchInterceptor
func (d *DisgoImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) {
	d.voiceDispatchInterceptor = voiceDispatchInterceptor
}

func (d *DisgoImpl) HTTPClient() rest.Client {
	return d.httpClient
}
func (d *DisgoImpl) RestServices() rest.Services {
	return d.restServices
}

// Gateway returns the websocket information
func (d *DisgoImpl) Gateway() gateway.Gateway {
	return d.gateway
}

// Connect opens the gateway connection to discord
func (d *DisgoImpl) Connect() error {
	return d.Gateway().Open()
}

// HasGateway returns whether api.Disgo has an active gateway.Gateway connection
func (d *DisgoImpl) HasGateway() bool {
	return d.gateway != nil
}

// HTTPServer returns the httpserver.Server
func (d *DisgoImpl) HTTPServer() httpserver.Server {
	return d.httpServer
}

// Start starts the interaction webhook server
func (d *DisgoImpl) Start() error {
	if d.HTTPServer() == nil {

	}
	d.HTTPServer().Start()
	return nil
}

// HasHTTPServer returns whether Disgo has an active httpserver.Server
func (d *DisgoImpl) HasHTTPServer() bool {
	return d.httpServer != nil
}

// Cache returns the entity api.Cache used by disgo
func (d *DisgoImpl) Cache() Cache {
	return d.cache
}

// EntityBuilder returns the api.EntityBuilder
func (d *DisgoImpl) EntityBuilder() EntityBuilder {
	return d.entityBuilder
}

// AudioController returns the api.AudioController which can be used to connect/reconnect/disconnect to/fom an api.VoiceChannel
func (d *DisgoImpl) AudioController() AudioController {
	return d.audioController
}

// GetCommand fetches a specific global discord.ApplicationCommand
func (d *DisgoImpl) GetCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().GetGlobalCommand(d.ApplicationID(), commandID, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetCommands fetches all global discord.ApplicationCommand(s)
func (d *DisgoImpl) GetCommands(opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().GetGlobalCommands(d.ApplicationID(), opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = d.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// CreateCommand creates a new global discord.ApplicationCommand
func (d *DisgoImpl) CreateCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().CreateGlobalCommand(d.ApplicationID(), commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// EditCommand edits a specific global discord.ApplicationCommand
func (d *DisgoImpl) EditCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().UpdateGlobalCommand(d.ApplicationID(), commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteCommand creates a new global discord.ApplicationCommand
func (d *DisgoImpl) DeleteCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return d.RestServices().ApplicationService().DeleteGlobalCommand(d.ApplicationID(), commandID, opts...)
}

// SetCommands overrides all global discord.ApplicationCommand(s)
func (d *DisgoImpl) SetCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().SetGlobalCommands(d.ApplicationID(), commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = d.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (d *DisgoImpl) GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().GetGuildCommand(d.ApplicationID(), guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (d *DisgoImpl) GetGuildCommands(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().GetGuildCommands(d.ApplicationID(), guildID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = d.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (d *DisgoImpl) CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().CreateGuildCommand(d.ApplicationID(), guildID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// EditGuildCommand edits a specific Guild discord.ApplicationCommand
func (d *DisgoImpl) EditGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().UpdateGuildCommand(d.ApplicationID(), guildID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (d *DisgoImpl) DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return d.RestServices().ApplicationService().DeleteGuildCommand(d.ApplicationID(), guildID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (d *DisgoImpl) SetGuildCommands(guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().SetGuildCommands(d.ApplicationID(), guildID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = d.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the api.GuildCommandPermissions for an all discord.ApplicationCommand(s) in an api.Guild
func (d *DisgoImpl) GetGuildCommandsPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().GetGuildCommandsPermissions(d.ApplicationID(), guildID, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = d.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the api.GuildCommandPermissions for a specific discord.ApplicationCommand in an api.Guild
func (d *DisgoImpl) GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	permissions, err := d.RestServices().ApplicationService().GetGuildCommandPermissions(d.ApplicationID(), guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// SetGuildCommandsPermissions sets the discord.GuildCommandPermissions for all discord.ApplicationCommand(s)
func (d *DisgoImpl) SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions []discord.GuildCommandPermissionsSet, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().SetGuildCommandsPermissions(d.ApplicationID(), guildID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = d.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific discord.ApplicationCommand
func (d *DisgoImpl) SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions []discord.CommandPermission, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().SetGuildCommandPermissions(d.ApplicationID(), guildID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommandPermissions(*perms, CacheStrategyNoWs), nil
}

// GetTemplate gets an api.GuildTemplate by its code
func (d *DisgoImpl) GetTemplate(code string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := d.RestServices().GuildTemplateService().GetGuildTemplate(code, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// CreateGuildFromTemplate creates an api.Guild using an api.Template code
func (d *DisgoImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, rest.Error) {
	guild, err := d.RestServices().GuildTemplateService().CreateGuildFromTemplate(templateCode, createGuildFromTemplate, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

func (d *DisgoImpl) GetInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error) {
	invite, err := d.RestServices().InviteService().GetInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (d *DisgoImpl) DeleteInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error) {
	invite, err := d.RestServices().InviteService().DeleteInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}

func buildDisgoImpl(config DisgoConfig) (Disgo, error) {
	disgo := &DisgoImpl{}

	if config.Token == "" {
		return nil, discord.ErrNoBotToken
	}
	disgo.token = config.Token

	id, err := IDFromToken(disgo.token)
	if err != nil {
		disgo.Logger().Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}
	// TODO: figure out how we handle different application & client ids
	disgo.applicationID = *id
	disgo.clientID = *id

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	disgo.logger = config.Logger

	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.RateLimiter == nil {
		config.RateLimiter = rate.NewLimiter(config.RateLimiterConfig)
	}

	if config.RestClientConfig == nil {
		config.RestClientConfig = &rest.DefaultConfig
	}

	if config.RestClientConfig.Headers == nil {
		config.RestClientConfig.Headers = http.Header{}
	}

	if _, ok := config.RestClientConfig.Headers["authorization"]; !ok {
		config.RestClientConfig.Headers["authorization"] = []string{fmt.Sprintf("Bot %s", config.Token)}
	}

	if config.RestClient == nil {
		config.RestClient = rest.NewClient(config.RestClientConfig)
	}

	if config.RestServices == nil {
		config.RestServices = rest.NewServices(disgo.logger, config.RestClient)
	}
	disgo.restServices = config.RestServices

	if config.EventManager == nil {
		config.EventManager = NewEventManager(disgo, config.EventListeners)
	}
	disgo.eventManager = config.EventManager

	if config.Gateway == nil && config.GatewayConfig != nil {
		config.Gateway = gateway.New(config.Token, func(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader) {
			disgo.EventManager().HandleGateway(gatewayEventType, sequenceNumber, payload)
		}, config.GatewayConfig)
	}
	disgo.gateway = config.Gateway

	if config.HTTPServer == nil && config.HTTPServerConfig != nil {
		config.HTTPServer = httpserver.New(func(responseChannel chan discord.InteractionResponse, payload io.Reader) {
			disgo.EventManager().HandleHTTP(responseChannel, payload)
		}, config.HTTPServerConfig)
	}
	disgo.httpServer = config.HTTPServer

	if config.AudioController == nil {
		config.AudioController = NewAudioController(disgo)
	}
	disgo.audioController = config.AudioController

	if config.EntityBuilder == nil {
		config.EntityBuilder = NewEntityBuilder(disgo)
	}
	disgo.entityBuilder = config.EntityBuilder

	disgo.voiceDispatchInterceptor = config.VoiceDispatchInterceptor

	if config.CacheConfig == nil {
		config.CacheConfig = &CacheConfig{
			CacheFlags:         CacheFlagsDefault,
			MemberCachePolicy:  MemberCachePolicyDefault,
			MessageCachePolicy: MessageCachePolicyDefault,
		}
	}

	if config.Cache == nil {
		config.Cache = NewCache(disgo, *config.CacheConfig)
	}
	disgo.cache = config.Cache

	return disgo, nil
}
