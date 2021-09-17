package sharding

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

var _ ShardManager = (*shardManagerImpl)(nil)

func NewShardManager(config Config) ShardManager {
	return &shardManagerImpl{
		shards: map[int]gateway.Gateway{},
		config: config,
	}
}

type shardManagerImpl struct {
	shards map[int]gateway.Gateway
	config Config
}

func (m *shardManagerImpl) Close() {

}
func (m *shardManagerImpl) Open() error {

}

func (m *shardManagerImpl) StartShard(shardID int) error {

}
func (m *shardManagerImpl) StopShard(shardID int) {

}

func (m *shardManagerImpl) GetGuildShard(guildId discord.Snowflake) gateway.Gateway {

}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {

}
func (m *shardManagerImpl) Shards() []gateway.Gateway {

}
