package sharding

import (
	"context"
	"sync"

	srate2 "github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/disgo/internal/merrors"
	"github.com/DisgoOrg/snowflake"

	"github.com/DisgoOrg/disgo/gateway"
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
		config.RateLimiter = srate2.NewLimiter(&srate2.DefaultConfig)
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

func (m *shardManagerImpl) Config() Config {
	return m.config
}

func (m *shardManagerImpl) RateLimiter() srate2.Limiter {
	return m.config.RateLimiter
}

func (m *shardManagerImpl) Open(ctx context.Context) error {
	m.Logger().Infof("opening %s shards...", m.config.Shards)
	var wg sync.WaitGroup
	var errs merrors.Error

	for shardInt := range m.config.Shards.Set {
		shardID := shardInt
		if m.shards.Has(shardID) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer m.RateLimiter().UnlockBucket(shardID)
			if err := m.RateLimiter().WaitBucket(ctx, shardID); err != nil {
				errs.Add(err)
				return
			}

			shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
			m.shards.Set(shardID, shard)
			if err := shard.Open(ctx); err != nil {
				errs.Add(err)
			}
		}()
	}
	wg.Wait()
	return errs
}

func (m *shardManagerImpl) ReOpen(ctx context.Context) error {
	m.Logger().Infof("reopening %s shards...", m.config.Shards)
	var wg sync.WaitGroup
	var errs merrors.Error

	for shardID := range m.shards.Shards {
		wg.Add(1)
		shard := m.shards.Get(shardID)
		go func() {
			defer wg.Done()
			if shard != nil {
				if err := shard.Close(ctx); err != nil {
					errs.Add(err)
				}
			}
			if err := shard.Open(ctx); err != nil {
				errs.Add(err)
			}
		}()
	}
	wg.Wait()
	return errs
}

func (m *shardManagerImpl) Close(ctx context.Context) error {
	m.Logger().Infof("closing %v shards...", m.config.Shards)
	var wg sync.WaitGroup
	var errs merrors.Error

	for shardID := range m.shards.Shards {
		m.config.Shards.Delete(shardID)
		shard := m.shards.Delete(shardID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := shard.Close(ctx); err != nil {
				errs.Add(err)
			}
		}()
	}
	wg.Wait()
	return errs
}

func (m *shardManagerImpl) OpenShard(ctx context.Context, shardID int) error {
	m.Logger().Infof("opening shard %d...", shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.eventHandlerFunc, m.config.GatewayConfig)
	m.config.Shards.Add(shardID)
	m.shards.Set(shardID, shard)
	return shard.Open(ctx)
}

func (m *shardManagerImpl) ReOpenShard(ctx context.Context, shardID int) error {
	m.Logger().Infof("reopening shard %d...", shardID)
	shard := m.shards.Get(shardID)
	if shard != nil {
		if err := shard.Close(ctx); err != nil {
			return err
		}
	}
	return shard.Open(ctx)
}

func (m *shardManagerImpl) CloseShard(ctx context.Context, shardID int) error {
	m.Logger().Infof("closing shard %d...", shardID)
	m.config.Shards.Delete(shardID)
	shard := m.shards.Delete(shardID)
	if shard != nil {
		return shard.Close(ctx)
	}
	return nil
}

func (m *shardManagerImpl) GetGuildShard(guildId snowflake.Snowflake) gateway.Gateway {
	return m.Shard(ShardIDByGuild(guildId, m.config.ShardCount))
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	return m.shards.Get(shardID)
}

func (m *shardManagerImpl) Shards() *ShardsMap {
	return m.shards
}
