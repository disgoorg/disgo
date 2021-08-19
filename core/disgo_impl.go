package core

import (
	"context"

	"github.com/DisgoOrg/disgo/httpserver"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

var _ Disgo = (*DisgoImpl)(nil)

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
func (d *DisgoImpl) Start() {
	d.HTTPServer().Start()
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
func (d *DisgoImpl) GetCommand(ctx context.Context, commandID discord.Snowflake) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().GetGlobalCommand(ctx, d.ApplicationID(), commandID)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetCommands fetches all global discord.ApplicationCommand(s)
func (d *DisgoImpl) GetCommands(ctx context.Context, ) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().GetGlobalCommands(ctx, d.ApplicationID())
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
func (d *DisgoImpl) CreateCommand(ctx context.Context, commandCreate discord.ApplicationCommandCreate) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().CreateGlobalCommand(ctx, d.ApplicationID(), commandCreate)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// EditCommand edits a specific global discord.ApplicationCommand
func (d *DisgoImpl) EditCommand(ctx context.Context, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().UpdateGlobalCommand(ctx, d.ApplicationID(), commandID, commandUpdate)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteCommand creates a new global discord.ApplicationCommand
func (d *DisgoImpl) DeleteCommand(ctx context.Context, commandID discord.Snowflake) rest.Error {
	return d.RestServices().ApplicationService().DeleteGlobalCommand(ctx, d.ApplicationID(), commandID)
}

// SetCommands overrides all global discord.ApplicationCommand(s)
func (d *DisgoImpl) SetCommands(ctx context.Context, commandCreates ...discord.ApplicationCommandCreate) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().SetGlobalCommands(ctx, d.ApplicationID(), commandCreates...)
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
func (d *DisgoImpl) GetGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().GetGuildCommand(ctx, d.ApplicationID(), guildID, commandID)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (d *DisgoImpl) GetGuildCommands(ctx context.Context, guildID discord.Snowflake) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().GetGuildCommands(ctx, d.ApplicationID(), guildID)
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
func (d *DisgoImpl) CreateGuildCommand(ctx context.Context, guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().CreateGuildCommand(ctx, d.ApplicationID(), guildID, commandCreate)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// EditGuildCommand edits a specific Guild discord.ApplicationCommand
func (d *DisgoImpl) EditGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (*ApplicationCommand, rest.Error) {
	command, err := d.RestServices().ApplicationService().UpdateGuildCommand(ctx, d.ApplicationID(), guildID, commandID, commandUpdate)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (d *DisgoImpl) DeleteGuildCommand(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) rest.Error {
	return d.RestServices().ApplicationService().DeleteGuildCommand(ctx, d.ApplicationID(), guildID, commandID)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (d *DisgoImpl) SetGuildCommands(ctx context.Context, guildID discord.Snowflake, commandCreates ...discord.ApplicationCommandCreate) ([]*ApplicationCommand, rest.Error) {
	cmds, err := d.RestServices().ApplicationService().SetGuildCommands(ctx, d.ApplicationID(), guildID, commandCreates...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = d.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the api.GuildCommandPermissions for a all discord.ApplicationCommand(s) in an api.Guild
func (d *DisgoImpl) GetGuildCommandsPermissions(ctx context.Context, guildID discord.Snowflake) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().GetGuildCommandsPermissions(ctx, d.ApplicationID(), guildID)
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
func (d *DisgoImpl) GetGuildCommandPermissions(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake) (*GuildCommandPermissions, rest.Error) {
	permissions, err := d.RestServices().ApplicationService().GetGuildCommandPermissions(ctx, d.ApplicationID(), guildID, commandID)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// SetGuildCommandsPermissions sets the discord.GuildCommandPermissions for all discord.ApplicationCommand(s)
func (d *DisgoImpl) SetGuildCommandsPermissions(ctx context.Context, guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().SetGuildCommandsPermissions(ctx, d.ApplicationID(), guildID, commandPermissions...)
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
func (d *DisgoImpl) SetGuildCommandPermissions(ctx context.Context, guildID discord.Snowflake, commandID discord.Snowflake, permissions ...discord.CommandPermission) (*GuildCommandPermissions, rest.Error) {
	perms, err := d.RestServices().ApplicationService().SetGuildCommandPermissions(ctx, d.ApplicationID(), guildID, commandID, permissions...)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateCommandPermissions(*perms, CacheStrategyNoWs), nil
}

// GetTemplate gets an api.GuildTemplate by its code
func (d *DisgoImpl) GetTemplate(ctx context.Context, code string) (*GuildTemplate, rest.Error) {
	guildTemplate, err := d.RestServices().GuildTemplateService().GetGuildTemplate(ctx, code)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// CreateGuildFromTemplate creates an api.Guild using an api.Template code
func (d *DisgoImpl) CreateGuildFromTemplate(ctx context.Context, templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*Guild, rest.Error) {
	guild, err := d.RestServices().GuildTemplateService().CreateGuildFromTemplate(ctx, templateCode, createGuildFromTemplate)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

func (d *DisgoImpl) GetInvite(ctx context.Context, inviteCode string) (*Invite, rest.Error) {
	invite, err := d.RestServices().InviteService().GetInvite(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (d *DisgoImpl) DeleteInvite(ctx context.Context, inviteCode string) (*Invite, rest.Error) {
	invite, err := d.RestServices().InviteService().DeleteInvite(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	return d.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}
