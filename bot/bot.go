package bot

import (
	"context"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake"
)

var _ Client = (*clientImpl)(nil)

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
	Rest() rest.Rest
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

	Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error
	Disconnect(ctx context.Context, guildID snowflake.Snowflake) error

	RequestMembers(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, userIDs ...snowflake.Snowflake) error
	RequestMembersWithQuery(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, query string, limit int) error

	SetPresence(ctx context.Context, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error
	SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error

	MemberChunkingManager() MemberChunkingManager

	StartHTTPServer() error
	HTTPServer() httpserver.Server
	HasHTTPServer() bool
}

// clientImpl is the main discord client
type clientImpl struct {
	token         string
	applicationID snowflake.Snowflake
	clientID      snowflake.Snowflake
	selfUser      *discord.OAuth2User

	logger log.Logger

	restServices rest.Rest

	eventManager EventManager

	shardManager sharding.ShardManager
	gateway      gateway.Gateway

	httpServer httpserver.Server

	caches cache.Caches

	memberChunkingManager MemberChunkingManager
}

func (c *clientImpl) Logger() log.Logger {
	return c.logger
}

// Close will clean up all disgo internals and close the discord connection safely
func (c *clientImpl) Close(ctx context.Context) {
	if c.restServices != nil {
		c.restServices.Close(ctx)
	}
	if c.gateway != nil {
		c.gateway.Close(ctx)
	}
	if c.shardManager != nil {
		c.shardManager.Close(ctx)
	}
	if c.httpServer != nil {
		c.httpServer.Close(ctx)
	}
}

func (c *clientImpl) Token() string {
	return c.token
}
func (c *clientImpl) ApplicationID() snowflake.Snowflake {
	return c.applicationID
}
func (c *clientImpl) ID() snowflake.Snowflake {
	return c.clientID
}
func (c *clientImpl) SelfUser() *discord.OAuth2User {
	return c.selfUser
}

func (c *clientImpl) SetSelfUser(user discord.OAuth2User) {
	c.selfUser = &user
}

// SelfMember returns a core.OAuth2User for the client, if available
func (c *clientImpl) SelfMember(guildID snowflake.Snowflake) *discord.Member {
	if member, ok := c.caches.Members().Get(guildID, c.clientID); ok {
		return &member
	}
	return nil
}

func (c *clientImpl) Caches() cache.Caches {
	return c.caches
}

func (c *clientImpl) Rest() rest.Rest {
	return c.restServices
}

func (c *clientImpl) HandleReadyEvent(event discord.GatewayEventReady) {
	c.applicationID = event.Application.ID
	c.clientID = event.User.ID
	c.selfUser = &event.User
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (c *clientImpl) AddEventListeners(listeners ...EventListener) {
	c.eventManager.AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (c *clientImpl) RemoveEventListeners(listeners ...EventListener) {
	c.eventManager.RemoveEventListeners(listeners...)
}

func (c *clientImpl) EventManager() EventManager {
	return c.eventManager
}

// ConnectGateway opens the gateway connection to discord
func (c *clientImpl) ConnectGateway(ctx context.Context) error {
	if c.gateway == nil {
		return discord.ErrNoGateway
	}
	return c.gateway.Open(ctx)
}

func (c *clientImpl) Gateway() gateway.Gateway {
	return c.gateway
}

// HasGateway returns whether this Client has an active gateway.Gateway connection
func (c *clientImpl) HasGateway() bool {
	return c.gateway != nil
}

// ConnectShardManager opens the gateway connection to discord
func (c *clientImpl) ConnectShardManager(ctx context.Context) error {
	if c.shardManager == nil {
		return discord.ErrNoShardManager
	}
	c.shardManager.Open(ctx)
	return nil
}

func (c *clientImpl) ShardManager() sharding.ShardManager {
	return c.shardManager
}

// HasShardManager returns whether this Client is sharded
func (c *clientImpl) HasShardManager() bool {
	return c.shardManager != nil
}

func (c *clientImpl) Shard(guildID snowflake.Snowflake) (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.gateway, nil
	} else if c.HasShardManager() {
		if shard := c.shardManager.GetGuildShard(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (c *clientImpl) Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.GatewayMessageDataVoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: &channelID,
	})
}

func (c *clientImpl) Disconnect(ctx context.Context, guildID snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.GatewayMessageDataVoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: nil,
	})
}

func (c *clientImpl) RequestMembers(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, userIDs ...snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.GatewayMessageDataRequestGuildMembers{
		GuildID:   guildID,
		Presences: presence,
		UserIDs:   userIDs,
		Nonce:     nonce,
	})
}
func (c *clientImpl) RequestMembersWithQuery(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, query string, limit int) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.GatewayMessageDataRequestGuildMembers{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presence,
		Nonce:     nonce,
	})
}

func (c *clientImpl) SetPresence(ctx context.Context, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	return c.gateway.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

// SetPresenceForShard sets the Presence of this Client for the provided shard
func (c *clientImpl) SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error {
	if !c.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := c.shardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

func (c *clientImpl) MemberChunkingManager() MemberChunkingManager {
	return c.memberChunkingManager
}

// StartHTTPServer starts the interaction webhook server
func (c *clientImpl) StartHTTPServer() error {
	if c.httpServer == nil {
		return discord.ErrNoHTTPServer
	}
	c.httpServer.Start()
	return nil
}

func (c *clientImpl) HTTPServer() httpserver.Server {
	return c.httpServer
}

// HasHTTPServer returns whether Client has an active httpserver.Server
func (c *clientImpl) HasHTTPServer() bool {
	return c.httpServer != nil
}
