package sharding

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type ShardManager interface {
	Close()
	Open() []error

	StartShard(shardID int) error
	StopShard(shardID int)

	GetGuildShard(guildId discord.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() []gateway.Gateway
}

func ShardForGuild(guildID discord.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}

func ShardRateLimitKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}
