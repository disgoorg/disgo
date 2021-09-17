package sharding

import (
	"io"
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

var _ ShardManager = (*shardManagerImpl)(nil)

func NewShardManager(token string, eventHandlerFunc gateway.EventHandlerFunc, config *Config) ShardManager {
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.Shards == nil || len(config.Shards) == 0{
		config.Shards = []int{0}
	}
	if config.ShardCount == 0 {
		config.ShardCount = len(config.Shards)
	}
	if config.GatewayCreateFunc == nil {
		config.GatewayCreateFunc = func() gateway.Gateway {
			return gateway.New(token, eventHandlerFunc, config.GatewayConfig)
		}
	}
	return &shardManagerImpl{
		shards: map[int]gateway.Gateway{},
		config: *config,
	}
}

type shardManagerImpl struct {
	shards map[int]gateway.Gateway
	config Config
}

func (m *shardManagerImpl) Close() {
	var wg sync.WaitGroup
	for i := range m.shards {
		shard := m.shards[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close()
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) Open() []error {
	var wg sync.WaitGroup
	var errs []error
	var mu
	for i := range m.config.Shards {
		if _, ok := m.shards[i]; ok {
			continue
		}
		shardID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gw := m.config.GatewayCreateFunc()
			m.shards[shardID] = gw
			err := gw.Open()
			if err != nil {
				errs = append(errs, err)
			}
		}()
	}
	wg.Wait()
	return errs
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
