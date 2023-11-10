package bot

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

var _ Client = (*clientImpl)(nil)

// Client is a high level interface for interacting with Discord's API.
// It combines the functionality of the rest, gateway/sharding, httpserver and cache into one easy to use interface.
// Create a new client with disgo.New.
type Client interface {
	// Logger returns the logger for the client.
	Logger() *slog.Logger

	// Close will clean up all disgo internals and close the discord gracefully.
	Close(ctx context.Context)

	// Token returns the configured bot token.
	Token() string

	// ApplicationID returns the application id.
	ApplicationID() snowflake.ID

	// ID returns the bot id.
	ID() snowflake.ID

	// Caches returns the cache.Caches used by the Client.
	Caches() cache.Caches

	// Rest returns the rest.Rest used by the Client.
	Rest() rest.Rest

	// AddEventListeners adds one or more EventListener(s) to the EventManager.
	AddEventListeners(listeners ...EventListener)

	// RemoveEventListeners removes one or more EventListener(s) from the EventManager
	RemoveEventListeners(listeners ...EventListener)

	// EventManager returns the EventManager used by the Client.
	EventManager() EventManager

	// VoiceManager returns the voice.Manager used by the Client.
	VoiceManager() voice.Manager

	// OpenGateway connects to the configured gateway.Gateway.
	OpenGateway(ctx context.Context) error

	// Gateway returns the gateway.Gateway used by the Client.
	Gateway() gateway.Gateway

	// HasGateway returns whether the Client has a configured gateway.Gateway.
	HasGateway() bool

	// OpenShardManager connects to the configured sharding.ShardManager.
	OpenShardManager(ctx context.Context) error

	// ShardManager returns the sharding.ShardManager used by the Client.
	ShardManager() sharding.ShardManager

	// HasShardManager returns whether the Client has a configured sharding.ShardManager.
	HasShardManager() bool

	// Shard returns the gateway.Gateway the specific guildID runs on.
	Shard(guildID snowflake.ID) (gateway.Gateway, error)

	// UpdateVoiceState sends a gateway.MessageDataVoiceStateUpdate to the specific gateway.Gateway.
	UpdateVoiceState(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error

	// RequestMembers sends a gateway.MessageDataRequestGuildMembers to the specific gateway.Gateway and requests the Member(s) of the specified guild.
	//  guildID  : is the snowflake of the guild to request the members of.
	//  presence : Whether to include discord.Presence data.
	//  nonce	 : The nonce to return to the discord.EventGuildMembersChunk.
	//  userIDs  : The snowflakes of the users to request the members of.
	RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error

	// RequestMembersWithQuery sends a gateway.MessageDataRequestGuildMembers to the specific gateway.Gateway and requests the Member(s) of the specified guild.
	//  guildID  : is the snowflake of the guild to request the members of.
	//  presence : Whether to include discord.Presence data.
	//  nonce    : The nonce to return to the discord.EventGuildMembersChunk.
	//  query    : The query to use for the request.
	//  limit    : The number of discord.Member(s) to return.
	RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error

	// SetPresence sends new presence data to the gateway.Gateway.
	SetPresence(ctx context.Context, opts ...gateway.PresenceOpt) error

	// SetPresenceForShard sends new presence data to the specific gateway.Gateway.
	SetPresenceForShard(ctx context.Context, shardId int, opts ...gateway.PresenceOpt) error

	// MemberChunkingManager returns the MemberChunkingManager used by the Client.
	MemberChunkingManager() MemberChunkingManager

	// OpenHTTPServer starts the configured HTTPServer used for interactions over webhooks.
	OpenHTTPServer() error

	// HTTPServer returns the configured HTTPServer used for interactions over webhooks.
	HTTPServer() httpserver.Server

	// HasHTTPServer returns whether the Client has a configured HTTPServer.
	HasHTTPServer() bool
}

type clientImpl struct {
	token         string
	applicationID snowflake.ID

	logger *slog.Logger

	restServices rest.Rest

	eventManager EventManager

	shardManager sharding.ShardManager
	gateway      gateway.Gateway

	httpServer httpserver.Server

	voiceManager voice.Manager

	caches cache.Caches

	memberChunkingManager MemberChunkingManager
}

func (c *clientImpl) Logger() *slog.Logger {
	return c.logger
}

func (c *clientImpl) Close(ctx context.Context) {
	if c.voiceManager != nil {
		c.voiceManager.Close(ctx)
	}
	if c.gateway != nil {
		c.gateway.Close(ctx)
	}
	if c.restServices != nil {
		c.restServices.Close(ctx)
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

func (c *clientImpl) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c *clientImpl) ID() snowflake.ID {
	if selfUser, ok := c.Caches().SelfUser(); ok {
		return selfUser.ID
	}
	return 0
}

func (c *clientImpl) Caches() cache.Caches {
	return c.caches
}

func (c *clientImpl) Rest() rest.Rest {
	return c.restServices
}

func (c *clientImpl) AddEventListeners(listeners ...EventListener) {
	c.eventManager.AddEventListeners(listeners...)
}

func (c *clientImpl) RemoveEventListeners(listeners ...EventListener) {
	c.eventManager.RemoveEventListeners(listeners...)
}

func (c *clientImpl) EventManager() EventManager {
	return c.eventManager
}

func (c *clientImpl) VoiceManager() voice.Manager {
	return c.voiceManager
}

func (c *clientImpl) OpenGateway(ctx context.Context) error {
	if c.gateway == nil {
		return discord.ErrNoGateway
	}
	return c.gateway.Open(ctx)
}

func (c *clientImpl) Gateway() gateway.Gateway {
	return c.gateway
}

func (c *clientImpl) HasGateway() bool {
	return c.gateway != nil
}

func (c *clientImpl) OpenShardManager(ctx context.Context) error {
	if c.shardManager == nil {
		return discord.ErrNoShardManager
	}
	c.shardManager.Open(ctx)
	return nil
}

func (c *clientImpl) ShardManager() sharding.ShardManager {
	return c.shardManager
}

func (c *clientImpl) HasShardManager() bool {
	return c.shardManager != nil
}

func (c *clientImpl) Shard(guildID snowflake.ID) (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.gateway, nil
	} else if c.HasShardManager() {
		if shard := c.shardManager.ShardByGuildID(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (c *clientImpl) UpdateVoiceState(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeVoiceStateUpdate, gateway.MessageDataVoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: channelID,
		SelfMute:  selfMute,
		SelfDeaf:  selfDeaf,
	})
}

func (c *clientImpl) RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeRequestGuildMembers, gateway.MessageDataRequestGuildMembers{
		GuildID:   guildID,
		Presences: presence,
		UserIDs:   userIDs,
		Nonce:     nonce,
	})
}

func (c *clientImpl) RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeRequestGuildMembers, gateway.MessageDataRequestGuildMembers{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presence,
		Nonce:     nonce,
	})
}

func (c *clientImpl) SetPresence(ctx context.Context, opts ...gateway.PresenceOpt) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	g := c.gateway
	return g.Send(ctx, gateway.OpcodePresenceUpdate, applyPresenceFromOpts(g, opts...))
}

func (c *clientImpl) SetPresenceForShard(ctx context.Context, shardId int, opts ...gateway.PresenceOpt) error {
	if !c.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := c.shardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, gateway.OpcodePresenceUpdate, applyPresenceFromOpts(shard, opts...))
}

func applyPresenceFromOpts(g gateway.Gateway, opts ...gateway.PresenceOpt) gateway.MessageDataPresenceUpdate {
	presenceUpdate := g.Presence()
	if presenceUpdate == nil {
		presenceUpdate = &gateway.MessageDataPresenceUpdate{}
	}
	for _, opt := range opts {
		opt(presenceUpdate)
	}
	return *presenceUpdate
}

func (c *clientImpl) MemberChunkingManager() MemberChunkingManager {
	return c.memberChunkingManager
}

func (c *clientImpl) OpenHTTPServer() error {
	if c.httpServer == nil {
		return discord.ErrNoHTTPServer
	}
	c.httpServer.Start()
	return nil
}

func (c *clientImpl) HTTPServer() httpserver.Server {
	return c.httpServer
}

func (c *clientImpl) HasHTTPServer() bool {
	return c.httpServer != nil
}
