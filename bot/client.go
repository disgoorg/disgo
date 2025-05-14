package bot

import (
	"context"
	"log/slog"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/disgo/voice"
)

// Client is a high level struct for interacting with Discord's API.
// It combines the functionality of the rest, gateway/sharding, httpserver and cache into one easy to use package.
// Create a new client with disgo.New.
type Client struct {
	Token                 string
	ApplicationID         snowflake.ID
	Logger                *slog.Logger
	Rest                  rest.Rest
	EventManager          EventManager
	ShardManager          sharding.ShardManager
	Gateway               gateway.Gateway
	HTTPServer            httpserver.Server
	VoiceManager          voice.Manager
	Caches                cache.Caches
	MemberChunkingManager MemberChunkingManager
}

func (c *Client) Close(ctx context.Context) {
	if c.VoiceManager != nil {
		c.VoiceManager.Close(ctx)
	}
	if c.Gateway != nil {
		c.Gateway.Close(ctx)
	}
	if c.Rest != nil {
		c.Rest.Close(ctx)
	}
	if c.ShardManager != nil {
		c.ShardManager.Close(ctx)
	}
	if c.HTTPServer != nil {
		c.HTTPServer.Close(ctx)
	}
}

func (c *Client) ID() snowflake.ID {
	if selfUser, ok := c.Caches.SelfUser(); ok {
		return selfUser.ID
	}
	return 0
}

func (c *Client) AddEventListeners(listeners ...EventListener) {
	c.EventManager.AddEventListeners(listeners...)
}

func (c *Client) RemoveEventListeners(listeners ...EventListener) {
	c.EventManager.RemoveEventListeners(listeners...)
}

func (c *Client) OpenGateway(ctx context.Context) error {
	if c.Gateway == nil {
		return discord.ErrNoGateway
	}
	return c.Gateway.Open(ctx)
}

func (c *Client) HasGateway() bool {
	return c.Gateway != nil
}

func (c *Client) OpenShardManager(ctx context.Context) error {
	if c.ShardManager == nil {
		return discord.ErrNoShardManager
	}
	c.ShardManager.Open(ctx)
	return nil
}

func (c *Client) HasShardManager() bool {
	return c.ShardManager != nil
}

func (c *Client) Shard(guildID snowflake.ID) (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.Gateway, nil
	} else if c.HasShardManager() {
		if shard := c.ShardManager.ShardByGuildID(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (c *Client) UpdateVoiceState(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error {
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

func (c *Client) RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error {
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

func (c *Client) RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error {
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

func (c *Client) RequestSoundboardSounds(ctx context.Context, guildIDs ...snowflake.ID) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	return c.Gateway.Send(ctx, gateway.OpcodeRequestSoundboardSounds, gateway.MessageDataRequestSoundboardSounds{
		GuildIDs: guildIDs,
	})
}

func (c *Client) SetPresence(ctx context.Context, opts ...gateway.PresenceOpt) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	g := c.Gateway
	return g.Send(ctx, gateway.OpcodePresenceUpdate, applyPresenceFromOpts(g, opts...))
}

func (c *Client) SetPresenceForShard(ctx context.Context, shardId int, opts ...gateway.PresenceOpt) error {
	if !c.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := c.ShardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, gateway.OpcodePresenceUpdate, applyPresenceFromOpts(shard, opts...))
}

func (c *Client) OpenHTTPServer() error {
	if c.HTTPServer == nil {
		return discord.ErrNoHTTPServer
	}
	c.HTTPServer.Start()
	return nil
}

func (c *Client) HasHTTPServer() bool {
	return c.HTTPServer != nil
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
