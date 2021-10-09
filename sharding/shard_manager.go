package sharding

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/sharding/srate"
	"github.com/DisgoOrg/log"
)

type ShardManager interface {
	Logger() log.Logger
	Config() Config
	RateLimiter() srate.Limiter

	Open() []error
	OpenCtx(ctx context.Context) []error
	Close()

	OpenShard(shardID int) error
	OpenShardCtx(ctx context.Context, shardID int) error

	ReopenShard(shardID int) error
	ReopenShardCtx(ctx context.Context, shardID int) error

	CloseShard(shardID int)

	GetGuildShard(guildId discord.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() *ShardsMap
}

func ShardIDByGuild(guildID discord.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}
