package sharding

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type ShardManager interface {
	Close()
	Open() []error

	OpenShard(shardID int) error
	ReopenShard(shardID int)
	CloseShard(shardID int)

	GetGuildShard(guildId discord.Snowflake) gateway.Gateway

	Shard(shardID int) gateway.Gateway
	Shards() map[int]gateway.Gateway
}

func ShardIDByGuild(guildID discord.Snowflake, shardCount int) int {
	return int((guildID.Int64() >> int64(22)) % int64(shardCount))
}
