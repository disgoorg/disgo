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

func New(token string, gatewayURL string, eventHandlerFunc gateway.EventHandlerFunc, config *Config) ShardManager {
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.Shards == nil || config.Shards.Len() == 0 {
		config.Shards = NewIntSet(0)
	}
	if config.ShardCount == 0 {
		config.ShardCount = config.Shards.Len()
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
		shards:           NewShardsMap(),
		token:            token,
		gatewayURL:       gatewayURL,
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards *ShardsMap

	token            string
	gatewayURL       string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) Logger() log.Logger {
	return m.config.Logger
}

func (m *shardManagerImpl) RateLimiter() rate.Limiter {
	return m.config.RateLimiter
}

func (m *shardManagerImpl) Open() []error {
	return m.OpenContext(context.Background())
}

func (m *shardManagerImpl) OpenContext(ctx context.Context) []error {
	m.Logger().Infof("opening %s shards...", m.config.Shards)
	var wg sync.WaitGroup
	var errs []error
	var errsMu sync.Mutex

	for shardInt := range m.config.Shards.Set {
		shardID := shardInt
		if m.shards.Has(shardID) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer m.RateLimiter().UnlockBucket(shardID)
			err := m.RateLimiter().WaitBucket(ctx, shardID)
			if err != nil {
				addErr(&errsMu, &errs, err)
				return
			}

			shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
			m.shards.Set(shardID, shard)
			err = shard.Open()
			if err != nil {
				addErr(&errsMu, &errs, err)
			}
		}()
	}
	wg.Wait()
	return errs
}

func addErr(mu *sync.Mutex, errs *[]error, err error) {
	mu.Lock()
	*errs = append(*errs, err)
	mu.Unlock()
}

func (m *shardManagerImpl) Close() {
	m.Logger().Infof("closing %v shards...", m.config.Shards)
	var wg sync.WaitGroup
	for shardID := range m.shards.Shards {
		m.config.Shards.Delete(shardID)
		shard := m.shards.Delete(shardID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close()
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) OpenShard(shardID int) error {
	return m.OpenShardContext(context.Background(), shardID)
}

func (m *shardManagerImpl) OpenShardContext(ctx context.Context, shardID int) error {
	m.Logger().Infof("opening shard %d...", shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
	m.config.Shards.Add(shardID)
	m.shards.Set(shardID, shard)
	return shard.OpenContext(ctx)
}

func (m *shardManagerImpl) ReopenShard(shardID int) error {
	return m.ReopenShardContext(context.Background(), shardID)
}

func (m *shardManagerImpl) ReopenShardContext(ctx context.Context, shardID int) error {
	m.Logger().Infof("reopening shard %d...", shardID)
	shard := m.shards.Get(shardID)
	if shard == nil {
		// TODO: should we start the shard if not already here?
		return nil
	}
	shard.Close()
	return shard.OpenContext(ctx)
}

func (m *shardManagerImpl) CloseShard(shardID int) {
	m.Logger().Infof("closing shard %d...", shardID)
	m.config.Shards.Delete(shardID)
	shard := m.shards.Get(shardID)
	if shard != nil {
		shard.Close()
	}
}

func (m *shardManagerImpl) GetGuildShard(guildId discord.Snowflake) gateway.Gateway {
	return m.Shard(ShardIDByGuild(guildId, m.config.ShardCount))
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	return m.shards.Get(shardID)
}

func (m *shardManagerImpl) Shards() *ShardsMap {
	return m.shards
}
