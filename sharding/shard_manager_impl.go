package sharding

import (
	"context"
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/sharding/rate"
	"github.com/DisgoOrg/log"
)

var _ ShardManager = (*shardManagerImpl)(nil)

func NewShardManager(token string, gatewayURLFunc func() string, eventHandlerFunc gateway.EventHandlerFunc, config *Config) ShardManager {
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.Shards == nil || len(config.Shards) == 0 {
		config.Shards = []int{0}
	}
	if config.ShardCount == 0 {
		config.ShardCount = len(config.Shards)
	}
	if config.GatewayConfig == nil {
		config.GatewayConfig = &gateway.DefaultConfig
	}
	if config.GatewayCreateFunc == nil {
		config.GatewayCreateFunc = func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
			return gateway.New(token, url, shardID, shardCount, eventHandlerFunc, config)
		}
	}
	if config.RateLimiter == nil {
		config.RateLimiter = rate.NewLimiter(&rate.DefaultConfig)
	}
	return &shardManagerImpl{
		shards:           map[int]gateway.Gateway{},
		token:            token,
		gatewayURL:       gatewayURLFunc(),
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards   map[int]gateway.Gateway
	shardsMu sync.RWMutex

	token            string
	gatewayURL       string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) Close() {
	var wg sync.WaitGroup
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
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
		m.shardsMu.RLock()
		if _, ok := m.shards[i]; ok {
			continue
		}
		m.shardsMu.RUnlock()
		shardID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer m.RateLimiter().UnlockBucket(shardID)
			err := m.RateLimiter().WaitBucket(context.Background(), shardID)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			gw := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
			m.shardsMu.Lock()
			m.shards[shardID] = gw
			m.shardsMu.Unlock()
			err = gw.Open()
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

func (m *shardManagerImpl) RateLimiter() rate.Limiter {
	return m.config.RateLimiter
}

func (m *shardManagerImpl) OpenShard(shardID int) error {
	gw := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
	m.shardsMu.Lock()
	m.shards[shardID] = gw
	m.shardsMu.Unlock()
	return gw.Open()
}

func (m *shardManagerImpl) ReopenShard(shardID int) error {
	return nil
}

func (m *shardManagerImpl) CloseShard(shardID int) {
	m.shardsMu.Lock()
	shard, ok := m.shards[shardID]
	delete(m.shards, shardID)
	m.shardsMu.Unlock()
	if ok {
		shard.Close()
	}
}

func (m *shardManagerImpl) GetGuildShard(guildId discord.Snowflake) gateway.Gateway {
	return m.Shard(ShardIDByGuild(guildId, m.config.ShardCount))
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	m.shardsMu.RLock()
	defer m.shardsMu.RUnlock()
	return m.shards[shardID]
}

func (m *shardManagerImpl) Shards() map[int]gateway.Gateway {
	return m.shards
}
