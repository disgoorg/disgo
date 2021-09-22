package sharding

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/sharding/rate"
)

type ShardManager interface {
	Close()
	Open() []error
	RateLimiter() rate.Limiter

	OpenShard(shardID int) error
	ReopenShard(shardID int) error
	CloseShard(shardID int)

	GetGuildShard(guildId discord.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() map[int]gateway.Gateway
}

func ShardIDByGuild(guildID discord.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}
