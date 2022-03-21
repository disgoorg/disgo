package bot

import (
	"context"

	"github.com/DisgoOrg/disgo/cache"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

var _ Client = (*ClientImpl)(nil)

type Client interface {
	Logger() log.Logger
	Close(ctx context.Context)

	Token() string
	ApplicationID() snowflake.Snowflake
	ID() snowflake.Snowflake
	SelfUser() *discord.OAuth2User
	SetSelfUser(user discord.OAuth2User)
	SelfMember(guildID snowflake.Snowflake) *discord.Member
	Caches() cache.Caches
	Rest() rest.Services
	HandleReadyEvent(event discord.GatewayEventReady)

	AddEventListeners(listeners ...EventListener)
	RemoveEventListeners(listeners ...EventListener)
	EventManager() EventManager

	ConnectGateway(ctx context.Context) error
	Gateway() gateway.Gateway
	HasGateway() bool

	ConnectShardManager(ctx context.Context) error
	ShardManager() sharding.ShardManager
	HasShardManager() bool
	Shard(guildID snowflake.Snowflake) (gateway.Gateway, error)

	AudioController() AudioController
	MemberChunkingManager() MemberChunkingManager

	SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error
	SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error

	StartHTTPServer() error
	HTTPServer() httpserver.Server
	HasHTTPServer() bool
}

// ClientImpl is the main discord client
type ClientImpl struct {
	BotToken         string
	BotApplicationID snowflake.Snowflake
	BotClientID      snowflake.Snowflake
	BotSelfUser      *discord.OAuth2User

	BotLogger log.Logger

	RestServices rest.Services

	BotEventManager EventManager

	BotShardManager sharding.ShardManager
	BotGateway      gateway.Gateway

	BotHTTPServer httpserver.Server

	BotCaches cache.Caches

	BotAudioController       AudioController
	BotMemberChunkingManager MemberChunkingManager
}

func (b *ClientImpl) Logger() log.Logger {
	return b.BotLogger
}

// Close will clean up all disgo internals and close the discord connection safely
func (b *ClientImpl) Close(ctx context.Context) {
	if b.RestServices != nil {
		b.RestServices.Close(ctx)
	}
	if b.BotGateway != nil {
		b.BotGateway.Close(ctx)
	}
	if b.BotShardManager != nil {
		b.BotShardManager.Close(ctx)
	}
	if b.BotHTTPServer != nil {
		b.BotHTTPServer.Close(ctx)
	}
}

func (b *ClientImpl) Token() string {
	return b.BotToken
}
func (b *ClientImpl) ApplicationID() snowflake.Snowflake {
	return b.BotApplicationID
}
func (b *ClientImpl) ID() snowflake.Snowflake {
	return b.BotClientID
}
func (b *ClientImpl) SelfUser() *discord.OAuth2User {
	return b.BotSelfUser
}

func (b *ClientImpl) SetSelfUser(user discord.OAuth2User) {
	b.BotSelfUser = &user
}

// SelfMember returns a core.OAuth2User for the client, if available
func (b *ClientImpl) SelfMember(guildID snowflake.Snowflake) *discord.Member {
	if member, ok := b.BotCaches.Members().Get(guildID, b.BotClientID); ok {
		return &member
	}
	return nil
}

func (b *ClientImpl) Caches() cache.Caches {
	return b.BotCaches
}

func (b *ClientImpl) Rest() rest.Services {
	return b.RestServices
}

func (b *ClientImpl) HandleReadyEvent(event discord.GatewayEventReady) {
	b.BotApplicationID = event.Application.ID
	b.BotClientID = event.User.ID
	b.BotSelfUser = &event.User
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (b *ClientImpl) AddEventListeners(listeners ...EventListener) {
	b.BotEventManager.AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (b *ClientImpl) RemoveEventListeners(listeners ...EventListener) {
	b.BotEventManager.RemoveEventListeners(listeners...)
}

func (b *ClientImpl) EventManager() EventManager {
	return b.BotEventManager
}

// ConnectGateway opens the BotGateway connection to discord
func (b *ClientImpl) ConnectGateway(ctx context.Context) error {
	if b.BotGateway == nil {
		return discord.ErrNoGateway
	}
	return b.BotGateway.Open(ctx)
}

func (b *ClientImpl) Gateway() gateway.Gateway {
	return b.BotGateway
}

// HasGateway returns whether this Client has an active gateway.Gateway connection
func (b *ClientImpl) HasGateway() bool {
	return b.BotGateway != nil
}

// ConnectShardManager opens the BotGateway connection to discord
func (b *ClientImpl) ConnectShardManager(ctx context.Context) error {
	if b.BotShardManager == nil {
		return discord.ErrNoShardManager
	}
	b.BotShardManager.Open(ctx)
	return nil
}

func (b *ClientImpl) ShardManager() sharding.ShardManager {
	return b.BotShardManager
}

// HasShardManager returns whether this Client is sharded
func (b *ClientImpl) HasShardManager() bool {
	return b.BotShardManager != nil
}

func (b *ClientImpl) Shard(guildID snowflake.Snowflake) (gateway.Gateway, error) {
	if b.HasGateway() {
		return b.BotGateway, nil
	} else if b.HasShardManager() {
		if shard := b.BotShardManager.GetGuildShard(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (b *ClientImpl) AudioController() AudioController {
	return b.BotAudioController
}

func (b *ClientImpl) MemberChunkingManager() MemberChunkingManager {
	return b.BotMemberChunkingManager
}

func (b *ClientImpl) SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !b.HasGateway() {
		return discord.ErrNoGateway
	}
	return b.BotGateway.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

// SetPresenceForShard sets the Presence of this Client for the provided shard
func (b *ClientImpl) SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !b.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := b.BotShardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodePresenceUpdate, presenceUpdate))
}

// StartHTTPServer starts the interaction webhook server
func (b *ClientImpl) StartHTTPServer() error {
	if b.BotHTTPServer == nil {
		return discord.ErrNoHTTPServer
	}
	b.BotHTTPServer.Start()
	return nil
}

func (b *ClientImpl) HTTPServer() httpserver.Server {
	return b.BotHTTPServer
}

// HasHTTPServer returns whether Client has an active httpserver.Server
func (b *ClientImpl) HasHTTPServer() bool {
	return b.BotHTTPServer != nil
}
