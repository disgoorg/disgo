package sharding

import "github.com/DisgoOrg/disgo/gateway"

func NewShardManager() ShardManager {

}

type shardManagerImpl struct {
	shards []gateway.Gateway
	config Config
}
