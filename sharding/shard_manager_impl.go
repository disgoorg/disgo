package sharding

import (
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
	if config.Shards == nil || len(config.Shards) == 0 {
		config.Shards = []int{0}
	}
	if config.ShardCount == 0 {
		config.ShardCount = len(config.Shards)
	}
	if config.GatewayCreateFunc == nil {
		config.GatewayCreateFunc = func(token string, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
			return gateway.New(token, eventHandlerFunc, config)
		}
	}
	return &shardManagerImpl{
		shards:           map[int]gateway.Gateway{},
		token:            token,
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards           map[int]gateway.Gateway
	token            string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) Close() {
	var wg sync.WaitGroup
	for i := range m.shards {
		shard := m.shards[i]
		delete(m.shards, i)
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
	var mu sync.Mutex
	for i := range m.config.Shards {
		if _, ok := m.shards[i]; ok {
			continue
		}
		shardID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gw := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.config.GatewayConfig)
			m.shards[shardID] = gw
			err := gw.Open()
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return errs
}

func (m *shardManagerImpl) OpenShard(shardID int) error {
	gw := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.config.GatewayConfig)
	m.shards[shardID] = gw
	return gw.Open()
}

func (m *shardManagerImpl) CloseShard(shardID int) {
	shard, ok := m.shards[shardID]
	delete(m.shards, shardID)
	if ok {
		shard.Close()
	}
}

func (m *shardManagerImpl) GetGuildShard(guildId discord.Snowflake) gateway.Gateway {
	shardID := ShardIDByGuild(guildId, m.config.ShardCount)
	return m.shards[shardID]
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	return m.shards[shardID]
}

func (m *shardManagerImpl) Shards() map[int]gateway.Gateway {
	return m.shards
}
