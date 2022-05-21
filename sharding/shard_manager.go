package sharding

import (
	"context"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type ShardManager interface {
	Logger() log.Logger

	Open(ctx context.Context)
	Close(ctx context.Context)

	OpenShard(ctx context.Context, shardID int) error
	CloseShard(ctx context.Context, shardID int)

	ShardByGuildID(guildId snowflake.ID) gateway.Gateway
	Shard(shardID int) gateway.Gateway
	Shards() map[int]gateway.Gateway
}

func ShardIDByGuild(guildID snowflake.ID, shardCount int) int {
	return int((uint64(guildID) >> 22) % uint64(shardCount))
}
