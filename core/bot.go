package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/internal/merrors"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

type Bot interface {
	Close(ctx context.Context) error

	ApplicationID() discord.Snowflake
	ClientID() discord.Snowflake
	SelfUser() SelfUser
	SelfMember(guildID discord.Snowflake) Member

	Logger() log.Logger
	RestServices() rest.Services
	EventManager() EventManager
	ShardManager() sharding.ShardManager
	Gateway() gateway.Gateway
	HTTPServer() httpserver.Server
	Caches() Caches
	EntityBuilder() EntityBuilder
	AudioController() AudioController
	MemberChunkingManager() MemberChunkingManager

	AddEventListeners(listeners ...EventListener)
	RemoveEventListeners(listeners ...EventListener)

	HasGateway() bool
	HasShardManager() bool

	ConnectGateway(ctx context.Context) error
	ConnectShardManager(ctx context.Context) error

	Shard(guildID discord.Snowflake) (gateway.Gateway, error)
	SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error
	SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error

	HasHTTPServer() bool
	StartHTTPServer() error

	GetCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error)
	GetCommands(opts ...rest.RequestOpt) ([]ApplicationCommand, error)
	CreateCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (ApplicationCommand, error)
	UpdateCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error)
	DeleteCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) error
	SetCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]ApplicationCommand, error)

	GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error)
	GetGuildCommands(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]ApplicationCommand, error)
	CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (ApplicationCommand, error)
	UpdateGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error)
	DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) error
	SetGuildCommands(guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]ApplicationCommand, error)

	GetGuildCommandsPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error)
	GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error)
	SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error)
	SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error)

	GetTemplate(code string, opts ...rest.RequestOpt) (*GuildTemplate, error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, error)

	GetInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, error)
	DeleteInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, error)

	GetNitroStickerPacks(opts ...rest.RequestOpt) ([]*StickerPack, error)
	GetSticker(stickerID discord.Snowflake, opts ...rest.RequestOpt) (*Sticker, error)

	CreateDMChannel(userID discord.Snowflake, opts ...rest.RequestOpt) (*DMChannel, error)
}

// BotImpl is the main discord client
type BotImpl struct {
	token         string
	applicationID discord.Snowflake
	clientID      discord.Snowflake
	selfUser      *SelfUser

	logger log.Logger

	restServices rest.Services

	eventManager EventManager

	shardManager sharding.ShardManager
	gateway      gateway.Gateway

	httpServer httpserver.Server

	caches Caches

	entityBuilder         EntityBuilder
	audioController       AudioController
	memberChunkingManager MemberChunkingManager
}

// Close will clean up all disgo internals and close the discord connection safely
func (b BotImpl) Close(ctx context.Context) error {
	var errs merrors.Error
	if b.RestServices() != nil {
		if err := b.RestServices().Close(ctx); err != nil {
			errs.Add(err)
		}
	}
	if b.Gateway() != nil {
		if err := b.Gateway().Close(ctx); err != nil {
			errs.Add(err)
		}
	}
	if b.ShardManager() != nil {
		if err := b.ShardManager().Close(ctx); err != nil {
			errs.Add(err)
		}
	}
	if b.HTTPServer() != nil {
		if err := b.HTTPServer().Close(ctx); err != nil {
			errs.Add(err)
		}
	}
	return nil
}

func (b *BotImpl) ApplicationID() discord.Snowflake {
	return b.applicationID
}
func (b *BotImpl) ClientID() discord.Snowflake {
	return b.clientID
}

func (b *BotImpl) SelfUser() SelfUser {
	return *b.selfUser
}

// SelfMember returns the Member for the specific Guild
func (b BotImpl) SelfMember(guildID discord.Snowflake) Member {
	member, _ := b.Caches().Members().Get(guildID, b.ClientID())
	return member
}

func (b *BotImpl) Logger() log.Logger {
	return b.logger
}

func (b *BotImpl) RestServices() rest.Services {
	return b.restServices
}

func (b *BotImpl) EventManager() EventManager {
	return b.eventManager
}

func (b *BotImpl) ShardManager() sharding.ShardManager {
	return b.shardManager
}

func (b *BotImpl) Gateway() gateway.Gateway {
	return b.gateway
}

func (b *BotImpl) HTTPServer() httpserver.Server {
	return b.httpServer
}

func (b *BotImpl) Caches() Caches {
	return b.caches
}

func (b *BotImpl) EntityBuilder() EntityBuilder {
	return b.entityBuilder
}

func (b *BotImpl) AudioController() AudioController {
	return b.audioController
}

func (b *BotImpl) MemberChunkingManager() MemberChunkingManager {
	return b.memberChunkingManager
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (b BotImpl) AddEventListeners(listeners ...EventListener) {
	b.EventManager().AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (b BotImpl) RemoveEventListeners(listeners ...EventListener) {
	b.EventManager().RemoveEventListeners(listeners...)
}

// ConnectGateway opens the gateway connection to discord
func (b BotImpl) ConnectGateway(ctx context.Context) error {
	if b.Gateway() == nil {
		return discord.ErrNoGateway
	}
	return b.Gateway().Open(ctx)
}

// ConnectShardManager opens the gateway connection to discord
func (b BotImpl) ConnectShardManager(ctx context.Context) error {
	if b.ShardManager() == nil {
		return discord.ErrNoShardManager
	}
	return b.ShardManager().Open(ctx)
}

// HasGateway returns whether this Bot has an active gateway.Gateway connection
func (b BotImpl) HasGateway() bool {
	return b.Gateway() != nil
}

// HasShardManager returns whether this Bot is sharded
func (b BotImpl) HasShardManager() bool {
	return b.ShardManager() != nil
}

func (b BotImpl) Shard(guildID discord.Snowflake) (gateway.Gateway, error) {
	if b.HasGateway() {
		return b.Gateway(), nil
	} else if b.HasShardManager() {
		shard := b.ShardManager().GetGuildShard(guildID)
		if shard == nil {
			return nil, discord.ErrShardNotFound
		}
		return shard, nil
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (b BotImpl) SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !b.HasGateway() {
		return discord.ErrNoGateway
	}
	return b.Gateway().Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

// SetPresenceForShard sets the Presence of this Bot for the provided shard
func (b BotImpl) SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !b.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := b.ShardManager().Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

// StartHTTPServer starts the interaction webhook server
func (b BotImpl) StartHTTPServer() error {
	if b.HTTPServer() == nil {
		return discord.ErrNoHTTPServer
	}
	b.HTTPServer().Start()
	return nil
}

// HasHTTPServer returns whether Bot has an active httpserver.Server
func (b BotImpl) HasHTTPServer() bool {
	return b.HTTPServer() != nil
}

// GetCommand fetches a specific global discord.ApplicationCommand
func (b BotImpl) GetCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().GetGlobalCommand(b.ApplicationID(), commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// GetCommands fetches all global discord.ApplicationCommand(s)
func (b BotImpl) GetCommands(opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := b.RestServices().ApplicationService().GetGlobalCommands(b.ApplicationID(), opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder().CreateApplicationCommand(command)
	}
	return commands, nil
}

// CreateCommand creates a new global discord.ApplicationCommand
func (b BotImpl) CreateCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().CreateGlobalCommand(b.ApplicationID(), commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// UpdateCommand updates a specific global discord.ApplicationCommand
func (b BotImpl) UpdateCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().UpdateGlobalCommand(b.ApplicationID(), commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// DeleteCommand creates a new global discord.ApplicationCommand
func (b BotImpl) DeleteCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) error {
	return b.RestServices().ApplicationService().DeleteGlobalCommand(b.ApplicationID(), commandID, opts...)
}

// SetCommands overrides all global discord.ApplicationCommand(s)
func (b BotImpl) SetCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := b.RestServices().ApplicationService().SetGlobalCommands(b.ApplicationID(), commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder().CreateApplicationCommand(command)
	}
	return commands, nil
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (b BotImpl) GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().GetGuildCommand(b.ApplicationID(), guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (b BotImpl) GetGuildCommands(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := b.RestServices().ApplicationService().GetGuildCommands(b.ApplicationID(), guildID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder().CreateApplicationCommand(command)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (b BotImpl) CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().CreateGuildCommand(b.ApplicationID(), guildID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// UpdateGuildCommand updates a specific Guild discord.ApplicationCommand
func (b BotImpl) UpdateGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := b.RestServices().ApplicationService().UpdateGuildCommand(b.ApplicationID(), guildID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommand(command), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (b BotImpl) DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) error {
	return b.RestServices().ApplicationService().DeleteGuildCommand(b.ApplicationID(), guildID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (b BotImpl) SetGuildCommands(guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := b.RestServices().ApplicationService().SetGuildCommands(b.ApplicationID(), guildID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder().CreateApplicationCommand(command)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the core.ApplicationCommandPermissions for an all discord.ApplicationCommand(s) in an core.Guild
func (b BotImpl) GetGuildCommandsPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error) {
	perms, err := b.RestServices().ApplicationService().GetGuildCommandsPermissions(b.ApplicationID(), guildID, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = b.EntityBuilder().CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand in a core.Guild
func (b BotImpl) GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := b.RestServices().ApplicationService().GetGuildCommandPermissions(b.ApplicationID(), guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommandPermissions(*permissions), nil
}

// SetGuildCommandsPermissions sets the discord.ApplicationCommandPermissions for all discord.ApplicationCommand(s)
func (b BotImpl) SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error) {
	perms, err := b.RestServices().ApplicationService().SetGuildCommandsPermissions(b.ApplicationID(), guildID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = b.EntityBuilder().CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand
func (b BotImpl) SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	perms, err := b.RestServices().ApplicationService().SetGuildCommandPermissions(b.ApplicationID(), guildID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateApplicationCommandPermissions(*perms), nil
}

// GetTemplate gets a core.GuildTemplate by its code
func (b BotImpl) GetTemplate(code string, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := b.RestServices().GuildTemplateService().GetGuildTemplate(code, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// CreateGuildFromTemplate creates a core.Guild using a core.GuildTemplate code
func (b BotImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, error) {
	guild, err := b.RestServices().GuildTemplateService().CreateGuildFromTemplate(templateCode, createGuildFromTemplate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

func (b BotImpl) GetInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, error) {
	invite, err := b.RestServices().InviteService().GetInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (b BotImpl) DeleteInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, error) {
	invite, err := b.RestServices().InviteService().DeleteInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (b BotImpl) GetNitroStickerPacks(opts ...rest.RequestOpt) ([]*StickerPack, error) {
	stickerPacks, err := b.RestServices().StickerService().GetNitroStickerPacks(opts...)
	if err != nil {
		return nil, err
	}
	coreStickerPacks := make([]*StickerPack, len(stickerPacks))
	for i, stickerPack := range stickerPacks {
		coreStickerPacks[i] = b.EntityBuilder().CreateStickerPack(stickerPack, CacheStrategyNoWs)
	}
	return coreStickerPacks, nil
}

func (b BotImpl) GetSticker(stickerID discord.Snowflake, opts ...rest.RequestOpt) (*Sticker, error) {
	sticker, err := b.RestServices().StickerService().GetSticker(stickerID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateSticker(*sticker, CacheStrategyNoWs), nil
}

func (b BotImpl) CreateDMChannel(userID discord.Snowflake, opts ...rest.RequestOpt) (*DMChannel, error) {
	sticker, err := b.RestServices().UserService().CreateDMChannel(userID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder().CreateChannel(*sticker, CacheStrategyNoWs).(*DMChannel), nil
}
