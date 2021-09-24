package sharding

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/sharding/rate"
	"github.com/DisgoOrg/log"
)

type ShardManager interface {
	Logger() log.Logger
	RateLimiter() rate.Limiter

	Open() []error
	OpenContext(ctx context.Context) []error
	Close()

	OpenShard(shardID int) error
	OpenShardContext(ctx context.Context, shardID int) error

	ReopenShard(shardID int) error
	ReopenShardContext(ctx context.Context, shardID int) error

	CloseShard(shardID int)

	GetGuildShard(guildId discord.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() *ShardsMap
}

func ShardIDByGuild(guildID discord.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}
