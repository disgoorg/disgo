package sharding

import (
	"context"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

type ShardManager interface {
	Logger() log.Logger
	Config() Config
	RateLimiter() srate.Limiter

	Open(ctx context.Context) error
	ReOpen(ctx context.Context) error
	Close(ctx context.Context)

	OpenShard(ctx context.Context, shardID int) error
	ReOpenShard(ctx context.Context, shardID int) error
	CloseShard(ctx context.Context, shardID int)

	GetGuildShard(guildId snowflake.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() *ShardsMap
}

func ShardIDByGuild(guildID snowflake.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}
