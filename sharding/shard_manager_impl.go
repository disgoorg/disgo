package sharding

import "github.com/DisgoOrg/disgo/gateway"

func NewShardManager() ShardManager {
	return nil
}

type shardManagerImpl struct {
	shards []gateway.Gateway
	config Config
}
