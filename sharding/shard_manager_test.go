package sharding

import (
	"fmt"
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

type shard struct {
	id    int
	count int
}

func TestSplit(t *testing.T) {
	var guildIDs []snowflake.ID
	for i := range 16 {
		guildIDs = append(guildIDs, snowflake.ID(i))
	}

	var shards = []shard{
		{0, 4},
		{1, 4},
		{2, 4},
		{3, 4},
	}

	printShards(shards, guildIDs)

	fmt.Println("-----")

	shards = []shard{
		{0, 12},
		{4, 12},
		{8, 12},
		{1, 4},
		{2, 4},
		{3, 4},
	}

	printShards(shards, guildIDs)
}

func printShards(shards []shard, guildIDs []snowflake.ID) {
	for _, s := range shards {
		for _, guildID := range guildIDs {
			shardID := ShardIDByGuild(guildID, s.count)
			if shardID == s.id {
				fmt.Printf("shard %d/%d handles guild %d\n", s.id, s.count, guildID)
			}
		}
	}
}
