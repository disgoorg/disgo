package core

import (
	"context"
	"io"
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	rrate "github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/sharding"
	srate "github.com/DisgoOrg/disgo/sharding/rate"
	"github.com/pkg/errors"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

var (
	defaultGatewayEventHandleFunc = func(bot *Bot) gateway.EventHandlerFunc {
		return func(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader) {
			bot.EventManager.HandleGateway(gatewayEventType, sequenceNumber, payload)
		}
	}
	defaultHTTPServerEventHandleFunc = func(bot *Bot) httpserver.EventHandlerFunc {
		return func(responseChannel chan<- discord.InteractionResponse, payload io.Reader) {
			bot.EventManager.HandleHTTP(responseChannel, payload)
		}
	}
)

func NewBot(token string, opts ...BotConfigOpt) (*Bot, error) {
	config := &BotConfig{}
	config.Apply(opts)

	return buildBot(token, *config)
}

// Bot is the main discord client
type Bot struct {
	Token         string
	ApplicationID discord.Snowflake
	ClientID      discord.Snowflake
	SelfUser      *SelfUser

	Logger log.Logger

	RestServices rest.Services

	EventManager EventManager

	ShardManager sharding.ShardManager
	Gateway      gateway.Gateway

	HTTPServer httpserver.Server

	Caches Caches

	EntityBuilder          EntityBuilder
	AudioController        AudioController
	MembersChunkingManager MembersChunkingManager
}

// Close will clean up all disgo internals and close the discord connection safely
func (b *Bot) Close() {
	if b.RestServices != nil {
		b.RestServices.Close()
	}
	if b.Gateway != nil {
		b.Gateway.Close()
	}
	if b.ShardManager != nil {
		b.ShardManager.Close()
	}
	if b.HTTPServer != nil {
		b.HTTPServer.Close()
	}
	if b.EventManager != nil {
		b.EventManager.Close()
	}
}

// SelfMember returns a core.OAuth2User for the client, if available
func (b *Bot) SelfMember(guildID discord.Snowflake) *Member {
	return b.Caches.MemberCache().Get(guildID, b.ClientID)
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (b *Bot) AddEventListeners(listeners ...EventListener) {
	b.EventManager.AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (b *Bot) RemoveEventListeners(listeners ...EventListener) {
	b.EventManager.RemoveEventListeners(listeners...)
}

// ConnectGateway opens the gateway connection to discord
func (b *Bot) ConnectGateway() error {
	return b.ConnectGatewayContext(context.Background())
}

func (b *Bot) ConnectGatewayContext(ctx context.Context) error {
	if b.Gateway == nil {
		return discord.ErrNoGateway
	}
	return b.Gateway.OpenContext(ctx)
}

// ConnectShardManager opens the gateway connection to discord
func (b *Bot) ConnectShardManager() []error {
	return b.ConnectShardManagerContext(context.Background())
}

func (b *Bot) ConnectShardManagerContext(ctx context.Context) []error {
	if b.ShardManager == nil {
		return []error{discord.ErrNoShardManager}
	}
	return b.ShardManager.OpenContext(ctx)
}

// HasGateway returns whether core.disgo has an active gateway.Gateway connection
func (b *Bot) HasGateway() bool {
	return b.Gateway != nil
}

func (b *Bot) HasShardManager() bool {
	return b.ShardManager != nil
}

func (b *Bot) SetPresence(presenceUpdate discord.PresenceUpdate) error {
	if !b.HasGateway() {
		return discord.ErrNoGateway
	}
	return b.Gateway.Send(discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

func (b *Bot) SetPresenceForShard(shardId int, presenceUpdate discord.PresenceUpdate) error {
	if !b.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := b.ShardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

// StartHTTPServer starts the interaction webhook server
func (b *Bot) StartHTTPServer() error {
	if b.HTTPServer == nil {
		return discord.ErrNoHTTPServer
	}
	b.HTTPServer.Start()
	return nil
}

// HasHTTPServer returns whether Bot has an active httpserver.Server
func (b *Bot) HasHTTPServer() bool {
	return b.HTTPServer != nil
}

// GetCommand fetches a specific global discord.ApplicationCommand
func (b *Bot) GetCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().GetGlobalCommand(b.ApplicationID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// GetCommands fetches all global discord.ApplicationCommand(s)
func (b *Bot) GetCommands(opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := b.RestServices.ApplicationService().GetGlobalCommands(b.ApplicationID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// CreateCommand creates a new global discord.ApplicationCommand
func (b *Bot) CreateCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().CreateGlobalCommand(b.ApplicationID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// EditCommand edits a specific global discord.ApplicationCommand
func (b *Bot) EditCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().UpdateGlobalCommand(b.ApplicationID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// DeleteCommand creates a new global discord.ApplicationCommand
func (b *Bot) DeleteCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return b.RestServices.ApplicationService().DeleteGlobalCommand(b.ApplicationID, commandID, opts...)
}

// SetCommands overrides all global discord.ApplicationCommand(s)
func (b *Bot) SetCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := b.RestServices.ApplicationService().SetGlobalCommands(b.ApplicationID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (b *Bot) GetGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().GetGuildCommand(b.ApplicationID, guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (b *Bot) GetGuildCommands(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := b.RestServices.ApplicationService().GetGuildCommands(b.ApplicationID, guildID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (b *Bot) CreateGuildCommand(guildID discord.Snowflake, commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().CreateGuildCommand(b.ApplicationID, guildID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// EditGuildCommand edits a specific Guild discord.ApplicationCommand
func (b *Bot) EditGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := b.RestServices.ApplicationService().UpdateGuildCommand(b.ApplicationID, guildID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommand(*command), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (b *Bot) DeleteGuildCommand(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return b.RestServices.ApplicationService().DeleteGuildCommand(b.ApplicationID, guildID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (b *Bot) SetGuildCommands(guildID discord.Snowflake, commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := b.RestServices.ApplicationService().SetGuildCommands(b.ApplicationID, guildID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = b.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the core.ApplicationCommandPermissions for an all discord.ApplicationCommand(s) in an core.Guild
func (b *Bot) GetGuildCommandsPermissions(guildID discord.Snowflake, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, rest.Error) {
	perms, err := b.RestServices.ApplicationService().GetGuildCommandsPermissions(b.ApplicationID, guildID, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = b.EntityBuilder.CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand in an core.Guild
func (b *Bot) GetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, rest.Error) {
	permissions, err := b.RestServices.ApplicationService().GetGuildCommandPermissions(b.ApplicationID, guildID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// SetGuildCommandsPermissions sets the discord.ApplicationCommandPermissions for all discord.ApplicationCommand(s)
func (b *Bot) SetGuildCommandsPermissions(guildID discord.Snowflake, commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, rest.Error) {
	perms, err := b.RestServices.ApplicationService().SetGuildCommandsPermissions(b.ApplicationID, guildID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = b.EntityBuilder.CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand
func (b *Bot) SetGuildCommandPermissions(guildID discord.Snowflake, commandID discord.Snowflake, permissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, rest.Error) {
	perms, err := b.RestServices.ApplicationService().SetGuildCommandPermissions(b.ApplicationID, guildID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateApplicationCommandPermissions(*perms), nil
}

// GetTemplate gets a core.GuildTemplate by its code
func (b *Bot) GetTemplate(code string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := b.RestServices.GuildTemplateService().GetGuildTemplate(code, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// CreateGuildFromTemplate creates a core.Guild using a core.GuildTemplate code
func (b *Bot) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...rest.RequestOpt) (*Guild, rest.Error) {
	guild, err := b.RestServices.GuildTemplateService().CreateGuildFromTemplate(templateCode, createGuildFromTemplate, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateGuild(*guild, CacheStrategyNoWs), nil
}

func (b *Bot) GetInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error) {
	invite, err := b.RestServices.InviteService().GetInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (b *Bot) DeleteInvite(inviteCode string, opts ...rest.RequestOpt) (*Invite, rest.Error) {
	invite, err := b.RestServices.InviteService().DeleteInvite(inviteCode, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateInvite(*invite, CacheStrategyNoWs), nil
}

func (b *Bot) GetNitroStickerPacks(opts ...rest.RequestOpt) ([]*StickerPack, rest.Error) {
	stickerPacks, err := b.RestServices.StickerService().GetNitroStickerPacks(opts...)
	if err != nil {
		return nil, err
	}
	coreStickerPacks := make([]*StickerPack, len(stickerPacks))
	for i, stickerPack := range stickerPacks {
		coreStickerPacks[i] = b.EntityBuilder.CreateStickerPack(stickerPack, CacheStrategyNoWs)
	}
	return coreStickerPacks, nil
}

func (b *Bot) GetSticker(stickerID discord.Snowflake, opts ...rest.RequestOpt) (*Sticker, rest.Error) {
	sticker, err := b.RestServices.StickerService().GetSticker(stickerID, opts...)
	if err != nil {
		return nil, err
	}
	return b.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

func buildBot(token string, config BotConfig) (*Bot, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := IDFromToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting application id from BotToken")
	}
	bot := &Bot{
		Token: token,
	}

	// TODO: figure out how we handle different application & client ids
	bot.ApplicationID = *id
	bot.ClientID = *id

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	bot.Logger = config.Logger

	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.RestClient == nil {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		if config.RestClientConfig.Logger == nil {
			config.RestClientConfig.Logger = config.Logger
		}
		if config.RestClientConfig.Headers == nil {
			config.RestClientConfig.Headers = http.Header{}
		}
		if _, ok := config.RestClientConfig.Headers["authorization"]; !ok {
			config.RestClientConfig.Headers["authorization"] = []string{discord.TokenTypeBot.Apply(token)}
		}

		if config.RestClientConfig.RateLimiterConfig == nil {
			config.RestClientConfig.RateLimiterConfig = &rrate.DefaultConfig
		}
		if config.RestClientConfig.RateLimiterConfig.Logger == nil {
			config.RestClientConfig.RateLimiterConfig.Logger = config.Logger
		}
		config.RestClient = rest.NewClient(config.RestClientConfig)
	}

	if config.RestServices == nil {
		config.RestServices = rest.NewServices(bot.Logger, config.RestClient)
	}
	bot.RestServices = config.RestServices

	if config.EventManager == nil {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &DefaultEventManagerConfig
		}
		config.EventManager = NewEventManager(bot, config.EventManagerConfig)
	}
	bot.EventManager = config.EventManager

	if config.Gateway == nil && config.GatewayConfig != nil {
		var gatewayRs *discord.Gateway
		gatewayRs, err = bot.RestServices.GatewayService().GetGateway()
		if err != nil {
			return nil, err
		}
		config.Gateway = gateway.New(token, gatewayRs.URL, 0, 0, defaultGatewayEventHandleFunc(bot), config.GatewayConfig)
	}
	bot.Gateway = config.Gateway

	if config.ShardManager == nil && config.ShardManagerConfig != nil {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = bot.RestServices.GatewayService().GetGatewayBot()
		if err != nil {
			return nil, err
		}

		if config.ShardManagerConfig.RateLimiterConfig == nil {
			config.ShardManagerConfig.RateLimiterConfig = &srate.DefaultConfig
		}
		if config.ShardManagerConfig.RateLimiterConfig.Logger == nil {
			config.ShardManagerConfig.RateLimiterConfig.Logger = config.Logger
		}
		if config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency == 0 {
			config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency = gatewayBotRs.SessionStartLimit.MaxConcurrency
		}

		// apply recommended shard count
		if !config.ShardManagerConfig.CustomShards {
			config.ShardManagerConfig.ShardCount = gatewayBotRs.Shards
			config.ShardManagerConfig.Shards = sharding.NewIntSet()
			for i := 0; i < gatewayBotRs.Shards; i++ {
				config.ShardManagerConfig.Shards.Add(i)
			}
		}
		if config.ShardManager == nil {
			config.ShardManager = sharding.New(token, gatewayBotRs.URL, defaultGatewayEventHandleFunc(bot), config.ShardManagerConfig)
		}
	}
	bot.ShardManager = config.ShardManager

	if config.HTTPServer == nil && config.HTTPServerConfig != nil {
		if config.HTTPServerConfig.Logger == nil {
			config.HTTPServerConfig.Logger = config.Logger
		}
		config.HTTPServer = httpserver.New(defaultHTTPServerEventHandleFunc(bot), config.HTTPServerConfig)
	}
	bot.HTTPServer = config.HTTPServer

	if config.AudioController == nil {
		config.AudioController = NewAudioController(bot)
	}
	bot.AudioController = config.AudioController

	if config.MembersChunkingManager == nil {
		config.MembersChunkingManager = NewMembersChunkingManager(bot)
	}
	bot.MembersChunkingManager = config.MembersChunkingManager

	if config.EntityBuilder == nil {
		config.EntityBuilder = NewEntityBuilder(bot)
	}
	bot.EntityBuilder = config.EntityBuilder

	if config.CacheConfig == nil {
		config.CacheConfig = &CacheConfig{
			CacheFlags:         CacheFlagsDefault,
			MemberCachePolicy:  MemberCachePolicyDefault,
			MessageCachePolicy: MessageCachePolicyDefault,
		}
	}

	if config.Caches == nil {
		config.Caches = NewCaches(*config.CacheConfig)
	}
	bot.Caches = config.Caches

	return bot, nil
}
