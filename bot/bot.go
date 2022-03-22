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

	Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error
	Disconnect(ctx context.Context, guildID snowflake.Snowflake) error

	RequestMembers(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, userIDs ...snowflake.Snowflake) error
	RequestMembersWithQuery(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, query string, limit int) error

	SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error
	SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error

	MemberChunkingManager() MemberChunkingManager

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

	BotMemberChunkingManager MemberChunkingManager
}

func (c *ClientImpl) Logger() log.Logger {
	return c.BotLogger
}

// Close will clean up all disgo internals and close the discord connection safely
func (c *ClientImpl) Close(ctx context.Context) {
	if c.RestServices != nil {
		c.RestServices.Close(ctx)
	}
	if c.BotGateway != nil {
		c.BotGateway.Close(ctx)
	}
	if c.BotShardManager != nil {
		c.BotShardManager.Close(ctx)
	}
	if c.BotHTTPServer != nil {
		c.BotHTTPServer.Close(ctx)
	}
}

func (c *ClientImpl) Token() string {
	return c.BotToken
}
func (c *ClientImpl) ApplicationID() snowflake.Snowflake {
	return c.BotApplicationID
}
func (c *ClientImpl) ID() snowflake.Snowflake {
	return c.BotClientID
}
func (c *ClientImpl) SelfUser() *discord.OAuth2User {
	return c.BotSelfUser
}

func (c *ClientImpl) SetSelfUser(user discord.OAuth2User) {
	c.BotSelfUser = &user
}

// SelfMember returns a core.OAuth2User for the client, if available
func (c *ClientImpl) SelfMember(guildID snowflake.Snowflake) *discord.Member {
	if member, ok := c.BotCaches.Members().Get(guildID, c.BotClientID); ok {
		return &member
	}
	return nil
}

func (c *ClientImpl) Caches() cache.Caches {
	return c.BotCaches
}

func (c *ClientImpl) Rest() rest.Services {
	return c.RestServices
}

func (c *ClientImpl) HandleReadyEvent(event discord.GatewayEventReady) {
	c.BotApplicationID = event.Application.ID
	c.BotClientID = event.User.ID
	c.BotSelfUser = &event.User
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (c *ClientImpl) AddEventListeners(listeners ...EventListener) {
	c.BotEventManager.AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (c *ClientImpl) RemoveEventListeners(listeners ...EventListener) {
	c.BotEventManager.RemoveEventListeners(listeners...)
}

func (c *ClientImpl) EventManager() EventManager {
	return c.BotEventManager
}

// ConnectGateway opens the BotGateway connection to discord
func (c *ClientImpl) ConnectGateway(ctx context.Context) error {
	if c.BotGateway == nil {
		return discord.ErrNoGateway
	}
	return c.BotGateway.Open(ctx)
}

func (c *ClientImpl) Gateway() gateway.Gateway {
	return c.BotGateway
}

// HasGateway returns whether this Client has an active gateway.Gateway connection
func (c *ClientImpl) HasGateway() bool {
	return c.BotGateway != nil
}

// ConnectShardManager opens the BotGateway connection to discord
func (c *ClientImpl) ConnectShardManager(ctx context.Context) error {
	if c.BotShardManager == nil {
		return discord.ErrNoShardManager
	}
	c.BotShardManager.Open(ctx)
	return nil
}

func (c *ClientImpl) ShardManager() sharding.ShardManager {
	return c.BotShardManager
}

// HasShardManager returns whether this Client is sharded
func (c *ClientImpl) HasShardManager() bool {
	return c.BotShardManager != nil
}

func (c *ClientImpl) Shard(guildID snowflake.Snowflake) (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.BotGateway, nil
	} else if c.HasShardManager() {
		if shard := c.BotShardManager.GetGuildShard(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (c *ClientImpl) Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: &channelID,
	})
}

func (c *ClientImpl) Disconnect(ctx context.Context, guildID snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: nil,
	})
}

func (c *ClientImpl) RequestMembers(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, userIDs ...snowflake.Snowflake) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.RequestGuildMembersCommandData{
		GuildID:   guildID,
		Presences: presence,
		UserIDs:   userIDs,
		Nonce:     nonce,
	})
}
func (c *ClientImpl) RequestMembersWithQuery(ctx context.Context, guildID snowflake.Snowflake, presence bool, nonce string, query string, limit int) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.RequestGuildMembersCommandData{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presence,
		Nonce:     nonce,
	})
}

func (c *ClientImpl) SetPresence(ctx context.Context, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	return c.BotGateway.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

// SetPresenceForShard sets the Presence of this Client for the provided shard
func (c *ClientImpl) SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.UpdatePresenceCommandData) error {
	if !c.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := c.BotShardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

func (c *ClientImpl) MemberChunkingManager() MemberChunkingManager {
	return c.BotMemberChunkingManager
}

// StartHTTPServer starts the interaction webhook server
func (c *ClientImpl) StartHTTPServer() error {
	if c.BotHTTPServer == nil {
		return discord.ErrNoHTTPServer
	}
	c.BotHTTPServer.Start()
	return nil
}

func (c *ClientImpl) HTTPServer() httpserver.Server {
	return c.BotHTTPServer
}

// HasHTTPServer returns whether Client has an active httpserver.Server
func (c *ClientImpl) HasHTTPServer() bool {
	return c.BotHTTPServer != nil
}
